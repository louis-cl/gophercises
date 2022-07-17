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
			"simple 1 link",
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
		{
			"2 links with nested elements",
			`
<html>
<head>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
	<h1>Social stuffs</h1>
	<div>
	<a href="https://www.twitter.com/joncalhoun">
		Check me out on twitter
		<i class="fa fa-twitter" aria-hidden="true"></i>
	</a>
	<a href="https://github.com/gophercises">
		Gophercises is on <strong>Github</strong>!
	</a>
	</div>
</body>
</html>
			`,
			[]Link{
				{Href: "https://www.twitter.com/joncalhoun", Text: "Check me out on twitter"},
				{Href: "https://github.com/gophercises", Text: "Gophercises is on Github!"},
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
