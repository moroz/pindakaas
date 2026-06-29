package services

import (
	"context"
	"database/sql"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
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
