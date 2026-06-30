package components

import (
	"fmt"

	"github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type IconProps struct {
	Name    string
	ViewBox string
	Classes string
}

func Icon(props *IconProps) Node {
	vb := "0 0 640 640"
	if props.ViewBox != "" {
		vb = props.ViewBox
	}

	return SVG(
		Class(twmerge.Merge("fill-current w-5 h-5", props.Classes)),
		Attr("viewBox", vb),
		El("use", Href(fmt.Sprintf("/assets/%s.svg#icon", props.Name))),
	)
}
