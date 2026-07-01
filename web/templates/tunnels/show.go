package tunnels

import (
	"fmt"
	"net"
	"strconv"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/components"
	"github.com/moroz/pindakaas/web/templates/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type fieldProps struct {
	Label      string
	Value      string
	Link       bool
	HelperText string
	Copy       bool
}

func field(props *fieldProps) Node {
	return Div(
		Class("flex flex-col gap-1 py-4 border-b border-slate-200 last:border-0"),
		Dt(
			Class("flex items-center text-sm gap-3 font-semibold uppercase tracking-wider text-slate-500"), Text(props.Label),
			If(
				props.Copy,
				components.CopyButton(props.Value),
			),
		),
		Dd(
			If(
				props.Link,
				A(
					Class("link"),
					Target("_blank"),
					Attr("rel", "noopener noreferrer"),
					Href(props.Value),
					Text(props.Value),
				),
			),
			If(!props.Link, Span(Class("font-mono text-slate-800"), Text(props.Value))),
			If(props.HelperText != "", P(
				Class("text-sm text-slate-500"),
				Text(props.HelperText),
			)),
		),
	)
}

func Show(ctx *types.RequestContext, data *queries.Tunnel) Node {
	fqdn := data.Subdomain + "." + config.BaseDomain
	if config.HTTPSPort != 443 {
		fqdn = net.JoinHostPort(fqdn, strconv.Itoa(int(config.HTTPSPort)))
	}
	fqdn = "https://" + fqdn
	credentials := data.Username + ":" + data.PasswordEncrypted.Plaintext()
	sshCmd := fmt.Sprintf("ssh -tt -R0:localhost:8080 -p %d %s@%s", config.SSHPort, credentials, config.BaseDomain)

	return layout.AppLayout(ctx, data.Subdomain,
		A(Href("/"), Class("mt-4 mb-2 inline-block link"), Text("<< Back to tunnels")),
		Header(
			Class("flex justify-between"),
			H2(Class("text-3xl font-bold mb-4 flex flex-col gap-2"), Text(data.Subdomain), Small(Class("font-semibold text-slate-500 text-base"), Text("Tunnel details"))),

			Form(
				Action(fmt.Sprintf("/tunnels/%s", data.ID)), Method("POST"),
				Input(Type("hidden"), Name("_method"), Value("DELETE")),
				Button(
					Type("submit"),
					Class("button secondary text-red-600 hover:text-red-700"),
					Attr("onclick", "return confirm('Delete this tunnel? This cannot be undone.')"),
					components.Icon(&components.IconProps{Name: "xmark", Classes: "w-6 h-6"}),
					Text("Delete tunnel"),
				),
			),
		),
		Div(
			Class("rounded-lg border border-slate-300 bg-white shadow-sm px-6"),
			El("dl",
				field(&fieldProps{
					Label: "URL",
					Value: fqdn,
					Copy:  true,
					Link:  true,
				}),
				field(&fieldProps{
					Label: "Username",
					Value: data.Username,
				}),
				field(&fieldProps{
					Label: "Password",
					Value: data.PasswordEncrypted.Plaintext(),
				}),
				field(&fieldProps{
					Label: "Credentials",
					Value: credentials,
					Copy:  true,
				}),
				field(&fieldProps{
					Label:      "SSH command",
					Value:      sshCmd,
					Copy:       true,
					HelperText: "Replace localhost:8080 with the actual URL of your local service if applicable.",
				}),
			),
		),
	)
}
