package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/types"
)

func SetRequestContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("context", types.NewRequestContext())
			return next(c)
		}
	}
}
