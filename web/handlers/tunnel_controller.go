package handlers

import (
	"database/sql"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/web/helpers"
	"github.com/moroz/pindakaas/web/templates/tunnels"
)

type tunnelController struct {
	db *sql.DB
}

func TunnelController(db *sql.DB) *tunnelController {
	return &tunnelController{db}
}

func (cc *tunnelController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	data, err := queries.New(cc.db).ListTunnels(c.Request().Context())
	if err != nil {
		return err
	}

	return tunnels.Index(ctx, &tunnels.IndexProps{
		Tunnels: data,
	}).Render(c.Response())
}
