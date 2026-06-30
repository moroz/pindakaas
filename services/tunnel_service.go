package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types"
)

type TunnelService struct {
	db       queries.DBTX
	registry types.TunnelRegistry
}

func NewTunnelService(db queries.DBTX, tunnelRegistry types.TunnelRegistry) *TunnelService {
	return &TunnelService{
		db:       db,
		registry: tunnelRegistry,
	}
}

func (s *TunnelService) AuthenticateHostBySSHUsername(ctx context.Context, userString string) (*queries.Tunnel, error) {
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

func randomHex4() (string, error) {
	var buf [4]byte
	_, err := rand.Read(buf[:])
	return hex.EncodeToString(buf[:]), err
}

func (s *TunnelService) CreateTunnelForUser(ctx context.Context, user *queries.User) (*types.TunnelCreateDTO, error) {
	subdomain, err := GenerateTunnelName()
	if err != nil {
		return nil, err
	}

	username, err := randomHex4()
	if err != nil {
		return nil, err
	}

	password, err := randomHex4()
	if err != nil {
		return nil, err
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	tunnel, err := queries.New(s.db).InsertTunnel(ctx, &queries.InsertTunnelParams{
		ID:           uuid.Must(uuid.NewV7()),
		Subdomain:    subdomain,
		Username:     username,
		PasswordHash: passwordHash,
		UserID:       user.ID,
	})
	if err != nil {
		return nil, err
	}

	return &types.TunnelCreateDTO{Tunnel: tunnel, PlaintextPassword: password}, nil
}

func (s *TunnelService) ListTunnelsForUser(ctx context.Context, user *queries.User) ([]*types.TunnelListDTO, error) {
	data, err := queries.New(s.db).ListTunnelsForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	var result []*types.TunnelListDTO

	for _, row := range data {
		result = append(result, &types.TunnelListDTO{
			Tunnel: row,
			Active: s.registry.GetTunnelStatus(row.Subdomain),
		})
	}

	return result, nil
}
