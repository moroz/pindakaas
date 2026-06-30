package layout

import (
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/templates/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AppHeader(ctx *types.RequestContext, title string, actions ...Node) Node {
	return Header(
		Class("h-16 bg-blue-600 text-white shadow fixed top-0 left-0 right-0"),
		Div(
			Class("container mx-auto flex justify-between items-center h-full"),
			H1(Class("text-3xl font-bold my-4"), Text("Pindakaas")),
			Div(
				Class("text-right ml-auto mr-6"),
				Span(
					Text(*ctx.User.GivenName+" "+*ctx.User.FamilyName),
					Br(),
					Text(ctx.User.Email),
				),
			),
			Form(
				Action("/sign-out"),
				Method("POST"),
				Input(Type("hidden"), Name("_method"), Value("DELETE")),
				Button(Class("button secondary gap-2"), Type("submit"),
					components.Icon(&components.IconProps{
						Name: "person-through-window",
					}),
					Text("Sign out"),
				),
			),
		),
	)

}
