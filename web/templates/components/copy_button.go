package components

import (
	"github.com/Oudwins/tailwind-merge-go/pkg/twmerge"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func CopyButton(text string, classes ...string) Node {
	return Div(Data("copy", text), Data("class", twmerge.Merge(classes...)))
}
