package handlers

import (
	"database/sql"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/helpers"
	"github.com/moroz/pindakaas/web/templates/tunnels"
)

type tunnelController struct {
	db            *sql.DB
	tunnelService *services.TunnelService
}

func TunnelController(db *sql.DB, tunnelRegistry types.TunnelRegistry) *tunnelController {
	return &tunnelController{
		db:            db,
		tunnelService: services.NewTunnelService(db, tunnelRegistry),
	}
}

func (cc *tunnelController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	data, err := cc.tunnelService.ListTunnelsForUser(c.Request().Context(), ctx.User)
	if err != nil {
		return err
	}

	return tunnels.Index(ctx, &tunnels.IndexProps{
		Tunnels: data,
	}).Render(c.Response())
}
