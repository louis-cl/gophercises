package parser

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestSimpleHtml(t *testing.T) {
	html := `
<html>
<body>
	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
</body>
</html>	
	`

	ans, err := AllLinksInHTML(strings.NewReader(html))
	if err != nil {
		t.Error("Got error", err)
	}
	expected := []Link{
		{Href: "/other-page", Text: "A link to another page"},
	}
	if !slices.Equal(ans, expected) {
		t.Errorf("got %v, want %v", ans, expected)
	}
}
