package registry

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/moroz/pindakaas/types"
)

type Registry struct {
	connections sync.Map
}

func New() *Registry {
	return &Registry{}
}

func (r *Registry) RegisterConnection(subdomain string, tunnel *types.Tunnel) (bool, error) {
	if _, ok := r.connections.Load(subdomain); ok {
		return false, fmt.Errorf("subdomain is already in use")
	}

	r.connections.Store(subdomain, tunnel)
	slog.Info("Registered connection", "subdomain", subdomain, "bind_addr", tunnel.BindAddr, "bind_port", tunnel.BindPort)

	return true, nil
}

func (r *Registry) DeregisterConnection(subdomain string) {
	r.connections.Delete(subdomain)
}

func (r *Registry) GetTunnelForSubdomain(subdomain string) (*types.Tunnel, bool) {
	tunnel, ok := r.connections.Load(subdomain)
	if !ok {
		return nil, false
	}
	return tunnel.(*types.Tunnel), ok
}
