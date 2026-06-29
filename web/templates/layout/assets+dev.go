//go:build !PROD

package layout

import (
	"github.com/moroz/pindakaas/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AssetEntryPoint(ctx *types.RequestContext) Node {
	entrypoint := "http://localhost:5173/src/main.ts"

	return Script(Type("module"), Src(entrypoint))
}
