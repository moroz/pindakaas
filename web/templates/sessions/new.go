package sessions

import (
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx *types.RequestContext) Node {
	return layout.RootLayout(ctx, "Sign in",
		Div(
			Class("pindakaas-bg min-h-screen"),
			H1(Text("Sign in")),

			A(
				Href("/oauth/google/redirect"),
				Text("Sign in with Google"),
			),
		))
}
