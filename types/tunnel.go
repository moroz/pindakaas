package types

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

type RequestPortForwardingPayload struct {
	BindAddr string
	BindPort uint32
}

type RequestPortForwardingSuccessPayload struct {
	BindPort uint32
}

type OpenForwardingChannelPayload struct {
	RemoteAddr string
	RemotePort uint32
	OriginAddr string
	OriginPort uint32
}

type Tunnel struct {
	Conn          *ssh.ServerConn
	BindAddr      string
	BindPort      uint32
	AllocatedPort uint32

	// logMu guards logSink, the destination for human-readable forwarding logs.
	// It is set when an interactive SSH session (e.g. `ssh -tt`) attaches and
	// cleared when that session ends. nil means no one is watching.
	logMu   sync.Mutex
	logSink chan<- string
}

var (
	_ http.RoundTripper = &Tunnel{}
)

// AttachLogSink directs forwarding logs to ch until DetachLogSink is called.
// The caller owns ch and is responsible for draining it.
func (t *Tunnel) AttachLogSink(ch chan<- string) {
	t.logMu.Lock()
	t.logSink = ch
	t.logMu.Unlock()
}

// DetachLogSink stops forwarding logs to ch, but only if it is still the active
// sink, so a stale session tearing down cannot detach a newer one.
func (t *Tunnel) DetachLogSink(ch chan<- string) {
	t.logMu.Lock()
	if t.logSink == ch {
		t.logSink = nil
	}
	t.logMu.Unlock()
}

// Logf sends a line to the attached log sink, if any. The send is non-blocking:
// when no session is watching, or the watcher's buffer is full, the line is
// dropped rather than slowing down request forwarding.
func (t *Tunnel) Logf(format string, args ...any) {
	t.logMu.Lock()
	defer t.logMu.Unlock()
	if t.logSink == nil {
		return
	}
	select {
	case t.logSink <- fmt.Sprintf(format, args...):
	default:
	}
}

func (t *Tunnel) OpenForwardingChannel(originAddr string, originPort uint32) (ForwardedConn, <-chan *ssh.Request, error) {
	payload := ssh.Marshal(OpenForwardingChannelPayload{
		RemoteAddr: t.BindAddr,
		RemotePort: t.AllocatedPort,
		OriginAddr: originAddr,
		OriginPort: originPort,
	})
	conn, reqs, err := t.Conn.OpenChannel("forwarded-tcpip", payload)
	return ForwardedConn{conn}, reqs, err
}

func parseRemoteAddrPort(addr string) (string, uint32, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	parsedPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return "", 0, err
	}

	return host, uint32(parsedPort), nil
}

type ForwardedResponse struct {
	io.ReadCloser
	channel ssh.Channel
}

func (r *ForwardedResponse) Close() error {
	err := r.ReadCloser.Close()
	r.channel.Close()
	return err
}

// forwardedRWC is an io.ReadWriteCloser over a forwarded SSH channel, used as
// the response body for "101 Switching Protocols" responses (e.g. WebSocket
// upgrades). httputil.ReverseProxy type-asserts the body to io.ReadWriteCloser
// to take over the connection and copy bytes bidirectionally.
//
// Reads go through r, the bufio.Reader that http.ReadResponse used: it may hold
// bytes read past the response headers (including the first frames of the
// upgraded protocol), so reading the raw channel directly would drop them.
type forwardedRWC struct {
	r       *bufio.Reader
	channel ssh.Channel
}

func (c *forwardedRWC) Read(p []byte) (int, error)  { return c.r.Read(p) }
func (c *forwardedRWC) Write(p []byte) (int, error) { return c.channel.Write(p) }
func (c *forwardedRWC) Close() error                { return c.channel.Close() }

func (t *Tunnel) RoundTrip(r *http.Request) (*http.Response, error) {
	host, port, err := parseRemoteAddrPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	channel, _, err := t.OpenForwardingChannel(host, port)
	if err != nil {
		return nil, err
	}

	err = r.Write(channel)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(channel)
	resp, err := http.ReadResponse(br, r)
	if err != nil {
		return nil, err
	}

	t.Logf("%s %s -> %d", r.Method, r.URL.RequestURI(), resp.StatusCode)

	if resp.StatusCode == http.StatusSwitchingProtocols {
		resp.Body = &forwardedRWC{r: br, channel: channel}
		return resp, nil
	}

	resp.Body = &ForwardedResponse{
		ReadCloser: resp.Body,
		channel:    channel,
	}

	return resp, nil
}
