package types

type TunnelRegistry interface {
	RegisterConnection(subdomain string, tunnel *Tunnel) (bool, error)
	DeregisterConnection(subdomain string)
	GetTunnelForSubdomain(subdomain string) (*Tunnel, bool)
}
