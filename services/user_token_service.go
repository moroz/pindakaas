package services

import (
	"context"
	"crypto/rand"
	"database/sql"

	"github.com/google/uuid"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types/dbtypes"
)

type UserTokenService struct {
	db *sql.DB
}

func NewUserTokenService(db *sql.DB) *UserTokenService {
	return &UserTokenService{db}
}

func generateToken() ([]byte, error) {
	var token = make([]byte, config.UserTokenLength)
	_, err := rand.Read(token)
	return token, err
}

func (s *UserTokenService) IssueAccessTokenForUser(ctx context.Context, user *queries.User) (*queries.UserToken, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertUserToken(ctx, &queries.InsertUserTokenParams{
		ID:      uuid.Must(uuid.NewV7()),
		UserID:  user.ID,
		Context: dbtypes.UserTokenContext_Access,
		Token:   token,
	})
}
