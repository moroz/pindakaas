package tunnels

import (
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type IndexProps struct {
	Tunnels []*queries.Tunnel
}

func Index(ctx *types.RequestContext, data *IndexProps) Node {
	return layout.RootLayout(ctx, "Tunnels", Div(
		H1(Class("text-3xl font-bold"), Text("Tunnels")),
	))
}
