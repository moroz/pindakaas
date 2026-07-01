package tunnels

import (
	"encoding/json"

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

type indexScriptProps struct {
	Tunnels    []*types.TunnelJSON `json:"tunnels"`
	BaseDomain string              `json:"baseDomain"`
}

func Index(ctx *types.RequestContext, data *IndexProps) Node {
	tunnelJSON := make([]*types.TunnelJSON, len(data.Tunnels))
	for i, t := range data.Tunnels {
		tunnelJSON[i] = t.ToJSON()
	}

	initialProps, _ := json.Marshal(&indexScriptProps{
		Tunnels:    tunnelJSON,
		BaseDomain: config.BaseDomain,
	})

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
		Div(ID("svelte-root")),
		Script(ID("index-table-props"), Type("application/json"), Raw(string(initialProps))),
	)
}
