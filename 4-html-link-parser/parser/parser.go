package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

func AllLinksInHTML(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var res []Link
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			res = append(res, linkInNode(node))
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return res, nil
}

func linkInNode(a *html.Node) Link {
	// extract href
	href := ""
	for _, attr := range a.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	// accumulate text
	var sb strings.Builder
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		} else if node.Type == html.ElementNode {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(a)
	return Link{Href: href, Text: sb.String()}
}
