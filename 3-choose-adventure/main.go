package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

const storyJsonPath = "gopher.json"
const htmlTemplatePath = "template.html"

func main() {
	story, err := readStory()
	if err != nil {
		panic(err)
	}

	templ := buildTemplate()

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		arc, err := arcFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid chapter", http.StatusBadRequest)
		}
		if arcStory, in := story[arc]; in {
			templ.Execute(w, arcStory)
			return
		} else {
			http.Error(w, "Chapter not found", http.StatusNotFound)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

// arcFromPath extracts which arc should be returned to a requested path
func arcFromPath(path string) (string, error) {
	if path == "/" || path == "" {
		return "intro", nil
	}
	if path[0] != '/' {
		return "", fmt.Errorf("invalid path %v", path)
	}
	return path[1:], nil
}

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

// readStory reads all arcs from storyJsonPath
// guaranteed to have one intro arc
func readStory() (map[string]storyArc, error) {
	bytes, err := os.ReadFile(storyJsonPath)
	if err != nil {
		return nil, err
	}

	arcs := make(map[string]storyArc)
	if err := json.Unmarshal(bytes, &arcs); err != nil {
		return nil, err
	}
	if _, ok := arcs["intro"]; !ok {
		return nil, errors.New("no intro arc in given json")
	}
	return arcs, nil
}

func buildTemplate() *template.Template {
	return template.Must(template.ParseFiles(htmlTemplatePath))
}
