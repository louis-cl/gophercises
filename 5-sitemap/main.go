package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"sitemap/parser"
)

func main() {
	var domain = flag.String("domain", "example.com", "domain to build a sitemap for")
	flag.Parse()

	links, err := linksIn(fmt.Sprintf("http://%s", *domain))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", links)
}

// linksIn returns the links present in the page at given url
// returned links are complete absolute http links
func linksIn(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	links, err := parser.AllLinksInHTML(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, link := range links {
		new_url, err := mergeUrl(link.Href, url)
		if err != nil {
			log.Println("Cannot merge url ", new_url, err)
		} else {
			result = append(result, new_url)
		}
	}
	return result, nil
}

func mergeUrl(child string, parent string) (string, error) {
	u, err := url.Parse(child)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" { // relative link
		return path.Join(parent, u.Path), nil
	} else if u.Scheme == "http" || u.Scheme == "https" { // absolute link
		return child, nil
	} else {
		return "", fmt.Errorf("unknown scheme %s", u.Scheme)
	}
}
