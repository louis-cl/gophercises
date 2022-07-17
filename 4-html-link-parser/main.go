package main

import (
	"fmt"
	"linkparser/parser"
	"os"
)

func main() {
	f, err := os.Open("./examples/ex1.html")
	if err != nil {
		panic(err)
	}
	links, err := parser.AllLinksInHTML(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", links)
}
