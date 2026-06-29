package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/helpers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type oidcController struct {
	db               *sql.DB
	config           *oauth2.Config
	userService      *services.UserService
	userTokenService *services.UserTokenService
}

func OIDCController(db *sql.DB) *oidcController {
	return &oidcController{
		db:               db,
		userService:      services.NewUserService(db),
		userTokenService: services.NewUserTokenService(db),
		config: &oauth2.Config{
			ClientID:     config.GoogleClientId,
			ClientSecret: config.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  config.PublicUrl + "/oauth/google/callback",
			Scopes:       []string{"email", "profile"},
		},
	}
}

func generateOAuthState() ([]byte, error) {
	var state = make([]byte, 4)
	_, err := rand.Read(state)
	return state, err
}

func (cc *oidcController) Redirect(c *echo.Context) error {
	if cc.config.ClientID == "" {
		log.Printf("Google Client ID is not set")
		return echo.NewHTTPError(500, "Client ID is not set")
	}

	ctx := helpers.GetRequestContext(c)

	state, err := generateOAuthState()
	if err != nil {
		return err
	}
	ctx.Session[config.OIDCStateSessionKey] = hex.EncodeToString(state)
	if err := ctx.SaveSession(c.Response()); err != nil {
		log.Printf("Error persisting session: %s", err)
		return err
	}

	url := cc.config.AuthCodeURL(hex.EncodeToString(state), oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

func decodeIDTokenClaims(token string) (*types.GoogleIDTokenClaims, error) {
	segs := strings.Split(token, ".")
	bytes, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		return nil, err
	}
	var claims types.GoogleIDTokenClaims
	if err := json.Unmarshal(bytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (cc *oidcController) Callback(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	state, _ := ctx.Session[config.OIDCStateSessionKey].(string)
	stateParam := c.QueryParam("state")

	if state != stateParam {
		log.Printf("Invalid OAuth2 state param in callback")
		return echo.NewHTTPError(400, "Invalid OAuth2 state param")
	}

	code := c.QueryParam("code")
	token, err := cc.config.Exchange(c.Request().Context(), code)
	if err != nil {
		log.Printf("Google token exchange returned error: %s", err)

		return echo.NewHTTPError(500, "Failed to fetch access token")
	}

	idToken, _ := token.Extra("id_token").(string)

	validator, _ := idtoken.NewValidator(c.Request().Context())
	_, err = validator.Validate(c.Request().Context(), idToken, cc.config.ClientID)
	if err != nil {
		log.Printf("ID token verification failed: %s", err)
		return err
	}

	claims, err := decodeIDTokenClaims(idToken)
	if err != nil {
		log.Printf("Failed to decode ID token: %s", err)
		return err
	}

	user, err := cc.userService.FindOrCreateUserByGoogleIDTokenClaims(c.Request().Context(), claims)
	if err != nil {
		log.Printf("Failed to create a user from claims: %s", err)
		return err
	}

	userToken, err := cc.userTokenService.IssueAccessTokenForUser(c.Request().Context(), user)
	if err != nil {
		log.Printf("Failed to issue access token for user: %s", err)
		return err
	}

	ctx.Session[config.AccessTokenSessionKey] = userToken.Token
	delete(ctx.Session, config.OIDCStateSessionKey)

	if err := ctx.SaveSession(c.Response()); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
