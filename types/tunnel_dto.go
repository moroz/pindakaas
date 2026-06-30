package types

import "github.com/moroz/pindakaas/db/queries"

type TunnelListDTO struct {
	*queries.Tunnel
	Active bool
}
