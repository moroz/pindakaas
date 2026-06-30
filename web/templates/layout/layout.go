package layout

import (
	"github.com/moroz/pindakaas/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RootLayout(ctx *types.RequestContext, title string, children ...Node) Node {
	pageTitle := "Pindakaas"
	if title != "" {
		pageTitle = title + " | " + pageTitle
	}

	return HTML(
		Lang("en"),
		Head(
			Meta(Charset("UTF-8")),
			TitleEl(Text(pageTitle)),
			AssetEntryPoint(ctx),
		),
		Body(Group(children)),
	)
}

func AppLayout(ctx *types.RequestContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		AppHeader(ctx),

		Div(
			Class("container mx-auto pt-16"),
			Group(children),
		))
}
