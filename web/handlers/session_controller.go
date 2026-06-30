package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/web/helpers"
	"github.com/moroz/pindakaas/web/templates/sessions"
)

type sessionController struct{}

func SessionController() *sessionController {
	return &sessionController{}
}

func (cc *sessionController) New(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	return sessions.New(ctx).Render(c.Response())
}
