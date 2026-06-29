package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/types/dbtypes"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) AuthenticateUserByAccessToken(ctx context.Context, token []byte) (*queries.User, error) {
	return queries.New(s.db).FindUserByUserToken(ctx, &queries.FindUserByUserTokenParams{
		Validity: int64(config.AccessTokenValidity.Seconds()),
		Token:    token,
		Context:  dbtypes.UserTokenContext_Access,
	})
}

func (s *UserService) FindOrCreateUserByGoogleIDTokenClaims(ctx context.Context, claims *types.GoogleIDTokenClaims) (*queries.User, error) {
	return queries.New(s.db).UpsertUser(ctx, &queries.UpsertUserParams{
		ID:         uuid.Must(uuid.NewV7()),
		Email:      claims.Email,
		GivenName:  &claims.GivenName,
		FamilyName: &claims.FamilyName,
		Avatar:     &claims.Avatar,
	})
}
