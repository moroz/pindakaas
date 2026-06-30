package tunnels

import (
	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/components"
	"github.com/moroz/pindakaas/web/templates/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.RequestContext, data *queries.Tunnel) Node {
	fqdn := "https://" + data.Subdomain + "." + config.BaseDomain
	return layout.AppLayout(ctx, "New Tunnel",
		A(Href("/"), Class("mt-4 inline-block underline text-primary"), Text("<< Back to tunnels")),
		H2(Class("text-2xl font-bold my-4"), Text("New Tunnel")),
		Table(
			Class("data-table w-full"),
			TBody(
				Tr(
					Th(Text("Subdomain")),
					Td(Class("font-mono"), Text(fqdn),
						components.CopyButton(fqdn),
					),
				),
				Tr(
					Th(Text("Username")),
					Td(Class("font-mono"), Text(data.Username)),
				),
				Tr(
					Th(Text("Password")),
					Td(Class("font-mono"), Text(data.PasswordEncrypted.Plaintext())),
				),
			),
		),
	)
}
