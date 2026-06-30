package components

import (
	"github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CopyButton(text string, classes ...string) Node {
	return Button(Attr("data-copy", text), Class(twmerge.Merge("button secondary ml-2 gap-1 font-sans h-8 px-2", twmerge.Merge(classes...))),
		Icon(&IconProps{Name: "copy"}),
		Text("Copy"),
	)
}
