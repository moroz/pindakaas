package helpers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/pindakaas/types"
)

func GetRequestContext(c *echo.Context) *types.RequestContext {
	return c.Get("context").(*types.RequestContext)
}
