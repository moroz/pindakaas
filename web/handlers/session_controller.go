package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/web/helpers"
	"github.com/moroz/pindakaas/web/templates/sessions"
)

type sessionController struct {
	service *services.UserTokenService
}

func SessionController(db *sql.DB) *sessionController {
	return &sessionController{
		service: services.NewUserTokenService(db),
	}
}

func (cc *sessionController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	return sessions.New(ctx).Render(c.Response())
}

func (cc *sessionController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	token, ok := ctx.Session[config.AccessTokenSessionKey].([]byte)
	if ok {
		cc.service.RevokeUserToken(c.Request().Context(), token)
	}

	delete(ctx.Session, config.AccessTokenSessionKey)
	if err := ctx.SaveSession(c.Response()); err != nil {
		log.Printf("Failed to save session cookie: %s", err)
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusFound, "/sign-in")
}
