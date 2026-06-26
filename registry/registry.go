package registry

import (
	"fmt"
	"log/slog"
	"sync"

	"golang.org/x/crypto/ssh"
)

type Registry struct {
	connections sync.Map
}

func New() *Registry {
	return &Registry{}
}

func (r *Registry) RegisterConnection(subdomain string, conn *ssh.ServerConn) (bool, error) {
	if _, ok := r.connections.Load(subdomain); ok {
		return false, fmt.Errorf("subdomain is already in use")
	}

	r.connections.Store(subdomain, conn)
	slog.Info("Registered connection", "subdomain", subdomain)

	return true, nil
}

func (r *Registry) DeregisterConnection(subdomain string) {
	r.connections.Delete(subdomain)
}

func (r *Registry) GetConnectionForSubdomain(subdomain string) (*ssh.ServerConn, bool) {
	conn, ok := r.connections.Load(subdomain)
	if !ok {
		return nil, false
	}
	return conn.(*ssh.ServerConn), ok
}
