package tunnels

import (
	"time"

	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/components"
	"github.com/moroz/pindakaas/web/templates/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type IndexProps struct {
	Tunnels []*types.TunnelListDTO
}

func statusBadge(active bool) Node {
	class := "badge inactive"
	icon := "bed"
	text := "Inactive"
	if active {
		class = "badge active"
		icon = "person-running"
		text = "Online"
	}
	return Span(
		Class(class),
		components.Icon(&components.IconProps{
			Name: icon,
		}),
		Text(text),
	)
}

func Index(ctx *types.RequestContext, data *IndexProps) Node {
	return layout.AppLayout(ctx, "Tunnels", Div(
		Data("hx-get", "/"),
		Data("hx-trigger", "every 10s"),
		Data("hx-select", ".index-table"),
		Data("hx-target", ".index-table"),
		Data("hx-swap", "outerHTML"),
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
				Map(data.Tunnels, func(tunnel *types.TunnelListDTO) Node {
					return Tr(
						Td(statusBadge(tunnel.Active)),
						Td(Class("font-mono"), Text(tunnel.Subdomain)),
						Td(Class("font-mono"), Text(tunnel.Username)),
						Td(Text(tunnel.InsertedAt.Format(time.RFC3339))),
					)
				}),
			),
		),
	))
}
