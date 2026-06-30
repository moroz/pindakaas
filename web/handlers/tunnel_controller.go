package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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

func (cc *tunnelController) Create(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	tunnel, err := cc.tunnelService.CreateTunnelForUser(c.Request().Context(), ctx.User)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/tunnels/%s", tunnel.ID))
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

func (cc *tunnelController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	id, err := uuid.Parse(c.Param("tunnel_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	tunnel, err := cc.tunnelService.GetTunnelForUser(c.Request().Context(), id, ctx.User)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}

	return tunnels.Show(ctx, tunnel).Render(c.Response())
}

func (cc *tunnelController) Delete(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	id, err := uuid.Parse(c.Param("tunnel_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := cc.tunnelService.DeleteTunnel(c.Request().Context(), id, ctx.User); err != nil {
		log.Printf("Failed to delete tunnel: %s", err)
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusFound, "/")
}
