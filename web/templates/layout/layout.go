package layout

import (
	"fmt"

	"github.com/moroz/pindakaas/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RootLayout(ctx *types.RequestContext, title string, children ...Node) Node {
	return HTML(
		Lang("en"),
		Head(
			Meta(Charset("UTF-8")),
			TitleEl(Text("Pindakaas")),
			AssetEntryPoint(ctx),
		),
		Body(
			Div(
				Class("container mx-auto"),
				Group(children),
			),
		),
	)
}

func AppLayout(ctx *types.RequestContext, title string, children ...Node) Node {
	return RootLayout(ctx, title, Div(
		Div(
			Class("container mx-auto"),
			Header(
				Class("flex justify-between items-center"),
				H1(Class("text-3xl font-bold my-4"), Text(title)),
				Div(
					Span(Text(fmt.Sprintf("%s %s (%s)", *ctx.User.GivenName, *ctx.User.FamilyName, ctx.User.Email))),
				),
			),
			Group(children),
		),
	))
}
