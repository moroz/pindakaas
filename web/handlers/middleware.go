package handlers

import (
	"database/sql"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/helpers"
	"github.com/moroz/pindakaas/web/sessions"
)

func SetRequestContext(store *sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("context", types.NewRequestContext(store))
			return next(c)
		}
	}
}

func FetchSessionFromCookies(store *sessions.Store, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)
			cookie, _ := c.Cookie(cookieName)
			payload, _ := store.DecodeSession(cookie)
			ctx.Session = payload
			return next(c)
		}
	}
}

func FetchUserFromSession(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := helpers.GetRequestContext(c)

			if token, ok := ctx.Session[config.AccessTokenSessionKey].([]byte); ok {
				if u, err := services.NewUserService(db).AuthenticateUserByAccessToken(c.Request().Context(), token); err == nil {
					ctx.User = u
				}
			}

			return next(c)
		}
	}
}

func CacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		path := c.Request().URL.Path

		// Cache versioned assets (containing hash in filename) for 1 year
		if strings.Contains(path, "-") && (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css")) {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Short cache for other assets
			c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		}

		return next(c)
	}
}
