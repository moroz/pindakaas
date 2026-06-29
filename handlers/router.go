package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Router(db *sql.DB) http.Handler {
	r := echo.New()

	r.Use(middleware.RequestLogger())
	r.Use(SetRequestContext())

	return r
}
