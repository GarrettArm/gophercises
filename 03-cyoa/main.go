package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Story map[string]Page

type Page struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func ParseJSON(filepath string) Story {
	var s Story
	b, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(err)
	}
	return s
}

func main() {
	tale := ParseJSON("gopher.json")

	fmt.Printf("%+v\n", tale)
}
