package types

type TunnelRegistry interface {
	GetTunnelStatus(subdomain string) bool
	RegisterConnection(subdomain string, tunnel *Tunnel)
	DeregisterConnection(subdomain string, tunnel *Tunnel)
	GetTunnelForSubdomain(subdomain string) (*Tunnel, bool)
}
