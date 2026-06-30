package services

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/bincyber/go-sqlcrypter"
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

	tunnel, err := queries.New(s.db).GetTunnelByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if match := subtle.ConstantTimeCompare([]byte(password), tunnel.PasswordEncrypted.Bytes()); match != 1 {
		return nil, ErrInvalidCredentials
	}

	return tunnel, nil
}

func randomHex(length int) (string, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return hex.EncodeToString(buf[:]), err
}

func (s *TunnelService) CreateTunnelForUser(ctx context.Context, user *queries.User) (*queries.Tunnel, error) {
	subdomain, err := GenerateTunnelName()
	if err != nil {
		return nil, err
	}

	username, err := randomHex(4)
	if err != nil {
		return nil, err
	}

	password, err := randomHex(8)
	if err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertTunnel(ctx, &queries.InsertTunnelParams{
		ID:                uuid.Must(uuid.NewV7()),
		Subdomain:         subdomain,
		Username:          username,
		PasswordEncrypted: sqlcrypter.NewEncryptedBytes(password),
		UserID:            user.ID,
	})
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

func (s *TunnelService) GetTunnelForUser(ctx context.Context, tunnelId uuid.UUID, user *queries.User) (*queries.Tunnel, error) {
	return queries.New(s.db).GetTunnelForUser(ctx, &queries.GetTunnelForUserParams{
		TunnelID: tunnelId,
		UserID:   user.ID,
	})
}
