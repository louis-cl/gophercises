package main

import (
	"encoding/json"
	"html/template"
	"os"
)

const storyJsonPath = "gopher.json"
const htmlTemplatePath = "template.html"

func main() {
	story, err := readStory()
	if err != nil {
		panic(err)
	}

	intro := story["intro"]
	templ := buildTemplate()
	templ.Execute(os.Stdout, intro)
}

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type allArcs map[string]storyArc

func readStory() (allArcs, error) {
	bytes, err := os.ReadFile(storyJsonPath)
	if err != nil {
		return nil, err
	}

	arcs := make(map[string]storyArc)
	if err := json.Unmarshal(bytes, &arcs); err != nil {
		return nil, err
	}
	return arcs, nil
}

func buildTemplate() *template.Template {
	return template.Must(template.ParseFiles(htmlTemplatePath))
}
