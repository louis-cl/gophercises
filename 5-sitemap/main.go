package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	var domain = flag.String("domain", "example.com", "domain to build a sitemap for")
	flag.Parse()

	resp, err := http.Get(fmt.Sprintf("http://%s", *domain))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
