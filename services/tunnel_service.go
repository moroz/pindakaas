package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

var ErrInvalidCredentials = errors.New("invalid username or password")

func (s *TunnelService) AuthenticateHostBySSHUsername(ctx context.Context, userString string) (*queries.Tunnel, error) {
	username, password, found := strings.Cut(userString, ":")
	if !found {
		return nil, ErrInvalidCredentials
	}

	host, err := queries.New(s.db).GetTunnelByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	match, err := argon2id.ComparePasswordAndHash(password, host.PasswordHash)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	return host, nil
}

func randomHex(length int) (string, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return hex.EncodeToString(buf[:]), err
}

func (s *TunnelService) CreateTunnelForUser(ctx context.Context, user *queries.User) (*types.TunnelDetailDTO, error) {
	subdomain, err := GenerateTunnelName()
	if err != nil {
		return nil, err
	}

	username, err := randomHex(4)
	if err != nil {
		return nil, err
	}

	password, err := randomHex(4)
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

	return &types.TunnelDetailDTO{
		Tunnel:            tunnel,
		PlaintextPassword: password,
	}, nil
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

func (s *TunnelService) DeleteTunnel(ctx context.Context, tunnelId uuid.UUID, user *queries.User) error {
	return queries.New(s.db).DeleteTunnelForUser(ctx, &queries.DeleteTunnelForUserParams{
		TunnelID: tunnelId,
		UserID:   user.ID,
	})
}

func (s *TunnelService) GetTunnelForUser(ctx context.Context, tunnelId uuid.UUID, user *queries.User) (*types.TunnelDetailDTO, error) {
	data, err := queries.New(s.db).GetTunnelForUser(ctx, &queries.GetTunnelForUserParams{
		TunnelID: tunnelId,
		UserID:   user.ID,
	})
	if err != nil {
		return nil, err
	}
	return &types.TunnelDetailDTO{
		Tunnel: data,
	}, nil
}
