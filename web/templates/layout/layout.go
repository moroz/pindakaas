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
		Body(
			Div(
				Class("container mx-auto"),
				Group(children),
			),
		),
	)
}
