package handlers

import (
	"database/sql"

	"github.com/labstack/echo/v5"
)

type tunnelController struct {
	db *sql.DB
}

func TunnelController(db *sql.DB) *tunnelController {
	return &tunnelController{db}
}

func (cc *tunnelController) Index(c *echo.Context) error {
	return nil
}
