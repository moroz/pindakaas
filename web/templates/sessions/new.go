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
			Class("pindakaas-bg min-h-screen flex items-center justify-center flex-col gap-4"),
			H1(Class("text-center text-white font-bold text-3xl"), Text("Pindakaas")),
			Div(Class("card"),
				H1(
					Class("text-2xl font-bold text-center mb-4"),
					Text("Sign in")),

				A(
					Class("button w-full"),
					Href("/oauth/google/redirect"),
					Text("Sign in with Google"),
				),
			),
		))
}
