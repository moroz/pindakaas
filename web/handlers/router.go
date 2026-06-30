package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/sessions"
)

type RouterProps struct {
	DB             *sql.DB
	Store          *sessions.Store
	TunnelRegistry types.TunnelRegistry
}

type Groupie interface {
	Group(string, ...echo.MiddlewareFunc) *echo.Group
}

func Group(r Groupie, prefix string, cb func(r *echo.Group)) {
	group := r.Group(prefix)
	cb(group)
}

func Router(props *RouterProps) http.Handler {
	r := echo.New()

	r.Use(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))
	r.Use(middleware.RequestID())
	r.Use(middleware.Recover())
	r.Use(middleware.RequestLogger())
	r.Use(SetRequestContext(props.Store))
	r.Use(FetchSessionFromCookies(props.Store, config.SessionCookieName))
	r.Use(FetchUserFromSession(props.DB))

	// Authenticated routes
	Group(r, "", func(r *echo.Group) {
		r.Use(RequireAuthenticatedUser)

		tunnels := TunnelController(props.DB)
		r.GET("/", tunnels.Index)
	})

	// Unauthenticated-only routes
	Group(r, "", func(r *echo.Group) {
		r.Use(RedirectToHomeIfAuthenticated)

		sessions := SessionController()
		r.GET("/sign-in", sessions.New)

		oauth2 := OIDCController(props.DB)
		r.GET("/oauth/google/redirect", oauth2.Redirect)
		r.GET("/oauth/google/callback", oauth2.Callback)
	})

	if config.IsProd {
		r.Static("/assets", "assets/dist/assets", CacheControlMiddleware)
	} else {
		r.Static("/assets", "assets/public/assets")
	}

	return r
}
