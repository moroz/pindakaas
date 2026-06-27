package types

import (
	"net/http"

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

func (t *Tunnel) RoundTrip(*http.Request) (*http.Response, error) {
	channel, err := t.OpenForwardingChannel()
}
