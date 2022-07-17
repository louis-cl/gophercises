package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	story, err := readStory()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", story)
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
	bytes, err := os.ReadFile("gopher.json")
	if err != nil {
		return nil, err
	}

	arcs := make(map[string]storyArc)
	if err := json.Unmarshal(bytes, &arcs); err != nil {
		return nil, err
	}
	return arcs, nil
}
