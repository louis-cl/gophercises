package scrapper

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sitemap/parser"
	"strings"
)

const depthToken string = "$"

func AllLinks(domain string, maxDepth int) []string {
	startingUrl := fmt.Sprintf("http://%s", domain)

	toProcess := list.New()
	processed := make(map[string]bool)

	toProcess.PushBack(startingUrl)
	toProcess.PushBack("$")
	processed[startingUrl] = false

	currentDepth := 0

	// token is after all elements of depth currentDepth
	for toProcess.Len() > 0 && currentDepth < maxDepth {
		item := toProcess.Front().Value.(string)
		toProcess.Remove(toProcess.Front())
		if item == depthToken {
			fmt.Println("reached depth", currentDepth)
			currentDepth++
			toProcess.PushBack(depthToken)
			continue
		}
		links, err := linksIn(item)
		// fmt.Println("got links for", item, links)
		if err != nil {
			log.Println("couldn't get links for", item, err)
		}
		processed[item] = true
		for _, link := range links {
			if _, in := processed[link]; !in && isLinkInDomain(link, domain) {
				processed[link] = false
				toProcess.PushBack(link)
			}
		}
	}

	result := make([]string, len(processed))
	i := 0
	for k := range processed {
		result[i] = k
		i++
	}
	return result
}

func isLinkInDomain(link string, domain string) bool {
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) < 2 {
		return u.Hostname() == domain
	}
	link_domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return link_domain == domain
}

// linksIn returns the links present in the page at given url
// returned links are complete absolute http links
func linksIn(web string) ([]string, error) {
	log.Println("scanning ", web)
	u, err := url.Parse(web)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(web)
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
		new_url, err := mergeUrl(link.Href, u)
		if err != nil {
			log.Println("Cannot merge url", new_url, err)
		} else {
			result = append(result, new_url)
		}
	}
	return result, nil
}

func mergeUrl(child string, parent *url.URL) (string, error) {
	u, err := url.Parse(child)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Scheme == "http" || u.Scheme == "https" {
		return parent.ResolveReference(u).String(), nil
	} else {
		return "", fmt.Errorf("unknown scheme %s", u.Scheme)
	}
}
