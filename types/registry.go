package types

type TunnelRegistry interface {
	RegisterConnection(subdomain string, tunnel *Tunnel)
	DeregisterConnection(subdomain string, tunnel *Tunnel)
	GetTunnelForSubdomain(subdomain string) (*Tunnel, bool)
}
