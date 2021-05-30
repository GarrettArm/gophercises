package main

import (
	"encoding/json"
	"fmt"
	"os"

	cyoa "github.com/GarrettArm/gophercises/03-cyoa"
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

func ParseJSON(filepath string) (Story, error) {
	var s Story
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func main() {
	cyoa.PrintTrash()
	tale, err := ParseJSON("gopher.json")
	if err != nil {
		fmt.Println("Error processing source json")
		panic(err)
	}

	fmt.Printf("%+v\n", tale)
}
