package registry

import (
	"log/slog"
	"sync"

	"github.com/moroz/pindakaas/types"
)

type Registry struct {
	connections sync.Map
}

var _ types.TunnelRegistry = &Registry{}

func New() *Registry {
	return &Registry{}
}

// RegisterConnection registers tunnel under subdomain using "last connection
// wins" semantics: if a tunnel is already registered for the subdomain, its SSH
// connection is closed and it is replaced by the new tunnel. The swap is atomic
// so concurrent registrations cannot lose an eviction.
func (r *Registry) RegisterConnection(subdomain string, tunnel *types.Tunnel) {
	previous, loaded := r.connections.Swap(subdomain, tunnel)
	slog.Info("Registered connection", "subdomain", subdomain, "bind_addr", tunnel.BindAddr, "bind_port", tunnel.BindPort)

	if loaded {
		old := previous.(*types.Tunnel)
		slog.Info("Evicting existing tunnel for subdomain", "subdomain", subdomain)
		old.Notify("This tunnel was taken over by a new connection from elsewhere. Disconnecting.")
		old.Conn.Close()
	}
}

// DeregisterConnection removes tunnel from subdomain only if it is still the
// registered tunnel. This guards against an evicted connection's teardown
// deleting the newer tunnel that replaced it.
func (r *Registry) DeregisterConnection(subdomain string, tunnel *types.Tunnel) {
	r.connections.CompareAndDelete(subdomain, tunnel)
}

func (r *Registry) GetTunnelForSubdomain(subdomain string) (*types.Tunnel, bool) {
	tunnel, ok := r.connections.Load(subdomain)
	if !ok {
		return nil, false
	}
	return tunnel.(*types.Tunnel), ok
}

func (r *Registry) GetTunnelStatus(subdomain string) bool {
	_, ok := r.connections.Load(subdomain)
	return ok
}
