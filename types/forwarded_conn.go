package types

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type ForwardedConn struct {
	ssh.Channel
}

var (
	_ net.Conn = &ForwardedConn{}
)

func (c ForwardedConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c ForwardedConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c ForwardedConn) SetDeadline(time.Time) error {
	return nil
}

func (c ForwardedConn) SetReadDeadline(time.Time) error {
	return nil
}

func (c ForwardedConn) SetWriteDeadline(time.Time) error {
	return nil
}
