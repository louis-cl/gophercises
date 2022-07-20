package main

import (
	"encoding/xml"
	"flag"
	"io"
	"os"
	"sitemap/scraper"
)

func main() {
	var domain = flag.String("domain", "example.com", "domain to build a sitemap for")
	flag.Parse()
	res := scraper.AllLinks(*domain)
	writeXml(os.Stdout, res)
}

func writeXml(out io.Writer, urls []string) {
	out.Write([]byte(xml.Header))
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	var out_urls []*xml_url
	for _, url := range urls {
		out_urls = append(out_urls, &xml_url{Loc: url})
	}
	enc.Encode(xml_urlset{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Urls: out_urls})
}

// <?xml version="1.0" encoding="UTF-8"?>
// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//   <url>
//     <loc>http://www.example.com/</loc>
//   </url>
//   <url>
//     <loc>http://www.example.com/dogs</loc>
//   </url>
// </urlset>

type xml_urlset struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	Urls    []*xml_url `xml:"url"`
}

type xml_url struct {
	Loc string `xml:"loc"`
}
