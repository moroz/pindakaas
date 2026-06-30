package tunnels

import (
	"fmt"
	"time"

	"github.com/moroz/pindakaas/config"
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
	return layout.AppLayout(ctx, "Tunnels",
		Header(
			Class("flex justify-between items-center"),
			H2(Class("text-2xl font-bold my-6"), Text("Tunnels")),
			Form(
				Action("/tunnels"), Method("POST"),
				Button(Type("submit"), Class("button gap-2"),
					components.Icon(&components.IconProps{
						Name: "plus",
					}),
					Text("New tunnel"),
				),
			),
		),
		Div(
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
						fqdn := "https://" + tunnel.Subdomain + "." + config.BaseDomain
						return Tr(
							Data("url", fmt.Sprintf("/tunnels/%s", tunnel.ID)),
							Td(statusBadge(tunnel.Active)),
							Td(
								Class("font-mono text-center"),
								Div(
									Class("inline-flex items-center"),
									Span(Attr("title", fqdn), Text(tunnel.Subdomain)),
									components.CopyButton(fqdn),
								),
							),
							Td(Class("font-mono"), Text(tunnel.Username)),
							Td(Text(tunnel.InsertedAt.Format(time.RFC3339))),
						)
					}),
				),
			),
		))
}
