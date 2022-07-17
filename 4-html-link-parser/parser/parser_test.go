package parser

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestTable(t *testing.T) {
	var tests = []struct {
		name string
		html string
		want []Link
	}{
		{
			"simple",
			`
<html>
<body>
	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
</body>
</html>
			`,
			[]Link{
				{Href: "/other-page", Text: "A link to another page"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := AllLinksInHTML(strings.NewReader(tt.html))
			if err != nil {
				t.Error("Got error", err)
			}
			if !slices.Equal(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
