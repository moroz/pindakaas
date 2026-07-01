package types

import (
	"github.com/google/uuid"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types/dbtypes"
)

type TunnelListDTO struct {
	*queries.Tunnel
	Active bool
}

type TunnelDetailDTO struct {
	*queries.Tunnel
	PlaintextPassword string
}

type TunnelJSON struct {
	ID         uuid.UUID             `json:"id"`
	Subdomain  string                `json:"subdomain"`
	Username   string                `json:"username"`
	InsertedAt dbtypes.UnixTimestamp `json:"insertedAt"`
	UpdatedAt  dbtypes.UnixTimestamp `json:"updatedAt"`
	UserID     uuid.UUID             `json:"userId"`
	Active     bool                  `json:"active"`
}

func (t *TunnelListDTO) ToJSON() *TunnelJSON {
	return &TunnelJSON{
		ID:         t.ID,
		Subdomain:  t.Subdomain,
		Username:   t.Username,
		InsertedAt: t.InsertedAt,
		UpdatedAt:  t.UpdatedAt,
		UserID:     t.UserID,
		Active:     t.Active,
	}
}
