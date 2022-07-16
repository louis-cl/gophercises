package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"urlshort/handler"
)

func main() {
	var mapFile = flag.String("routes", "", "path to yaml file with routes")
	flag.Parse()

	yamlConfig := readOrDefault(*mapFile)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := handler.YAMLHandler(yamlConfig, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using YAMLHandler as the fallback

	jsonConfig := []byte(`
	[{
		"path": "/json",
		"url": "https://jsoneditoronline.org/"
	},
		{
		"path": "/yaml",
		"url": "https://codebeautify.org/yaml-editor-online"
	}]`)
	jsonHandler, err := handler.JSONHandler(jsonConfig, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func readOrDefault(filePath string) []byte {
	if len(filePath) == 0 {
		return []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
