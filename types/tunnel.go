package types

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"

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
}

var (
	_ http.RoundTripper = &Tunnel{}
)

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

	resp, err := http.ReadResponse(bufio.NewReader(channel), r)
	if err != nil {
		return nil, err
	}

	resp.Body = &ForwardedResponse{
		ReadCloser: resp.Body,
		channel:    channel,
	}

	return resp, nil
}
