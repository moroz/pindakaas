package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/moroz/pindakaas/db/queries"
)

type HostService struct {
	db queries.DBTX
}

func NewHostService(db queries.DBTX) *HostService {
	return &HostService{db}
}

func (s *HostService) AuthenticateHostBySSHUsername(ctx context.Context, userString string) (*queries.Tunnel, error) {
	username, password, found := strings.Cut(userString, ":")
	if !found {
		return nil, fmt.Errorf("malformed username")
	}

	host, err := queries.New(s.db).GetTunnelByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, host.PasswordHash)
	if err != nil || !match {
		return nil, fmt.Errorf("invalid username or password")
	}

	return host, nil
}
