package main

import (
	"flag"
	"fmt"
	"sitemap/scraper"
)

func main() {
	var domain = flag.String("domain", "example.com", "domain to build a sitemap for")
	flag.Parse()
	res := scraper.AllLinks(*domain)
	for _, link := range res {
		fmt.Println(link)
	}
}
