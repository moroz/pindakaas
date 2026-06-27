package types

import "golang.org/x/crypto/ssh"

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

func (t *Tunnel) OpenForwardingChannel(originAddr string, originPort uint32) (ssh.Channel, <-chan *ssh.Request, error) {
	payload := ssh.Marshal(OpenForwardingChannelPayload{
		RemoteAddr: t.BindAddr,
		RemotePort: t.AllocatedPort,
		OriginAddr: originAddr,
		OriginPort: originPort,
	})
	return t.Conn.OpenChannel("forwarded-tcpip", payload)
}
