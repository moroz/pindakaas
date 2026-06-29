package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/web/helpers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oidcController struct {
	db     *sql.DB
	config *oauth2.Config
}

func OIDCController(db *sql.DB) *oidcController {
	return &oidcController{
		db: db,
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
