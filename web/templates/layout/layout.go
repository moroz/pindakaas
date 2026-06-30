package layout

import (
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
		Body(Group(children)),
	)
}

func AppLayout(ctx *types.RequestContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Header(
			Class("h-20 bg-blue-600 text-white shadow fixed top-0 left-0 right-0"),
			Div(
				Class("container mx-auto flex justify-between items-center h-full"),
				H1(Class("text-3xl font-bold my-4"), Text("Pindakaas")),
				Div(
					Class("text-right"),
					Span(
						Text(*ctx.User.GivenName+" "+*ctx.User.FamilyName),
						Br(),
						Text(ctx.User.Email),
					),
				),
			),
		),

		Div(
			Class("container mx-auto pt-24"),
			Header(
				Class("flex justify-between items-center"),
			),

			Group(children),
		))
}
