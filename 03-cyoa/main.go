package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

func parseFlags() (storyFile *string, port *int, template *string) {
	storyFile = flag.String("f", "gopher.json", "a json file containing a story")
	port = flag.Int("p", 3030, "port on which cyoa is served")
	template = flag.String("t", "story.tmpl", "template to be used")
	flag.Parse()
	return storyFile, port, template
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	if page, ok := h.s[path]; ok {
		err := h.t.Execute(w, page)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Page Not Found", http.StatusNotFound)
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func storyHandler(s Story, opts ...HandlerOption) http.Handler {
	t := template.Must(template.ParseFiles("story.tmpl"))
	h := handler{s, t}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func main() {
	storyFile, port, templateFile := parseFlags()

	story, err := ParseJSON(*storyFile)
	if err != nil {
		fmt.Println("Error processing source json")
		panic(err)
	}

	fmt.Printf("Serving on %d\n", *port)
	template := template.Must(template.ParseFiles(*templateFile))
	handler := storyHandler(story, WithTemplate(template))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
