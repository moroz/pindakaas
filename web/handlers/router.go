package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/web/sessions"
)

func Router(db *sql.DB, store *sessions.Store) http.Handler {
	r := echo.New()

	r.Use(middleware.RequestLogger())
	r.Use(SetRequestContext(store))
	r.Use(FetchSessionFromCookies(store, config.SessionCookieName))
	r.Use(FetchUserFromSession(db))

	oauth2 := OIDCController(db)
	r.GET("/oauth/google/redirect", oauth2.Redirect)

	return r
}
