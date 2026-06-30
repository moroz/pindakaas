package tunnels

import (
	"time"

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
	return layout.AppLayout(ctx, "Tunnels", Div(
		Table(
			Class("index-table w-full"),
			THead(
				Tr(
					Th(Text("Active")),
					Th(Text("Subdomain")),
					Th(Text("Username")),
					Th(Text("Created at")),
				),
			),
			TBody(
				Map(data.Tunnels, func(tunnel *queries.Tunnel) Node {
					return Tr(
						Td(),
						Td(Class("font-mono"), Text(tunnel.Subdomain)),
						Td(Class("font-mono"), Text(tunnel.Username)),
						Td(Text(tunnel.InsertedAt.Format(time.RFC3339))),
					)
				}),
			),
		),
	))
}
