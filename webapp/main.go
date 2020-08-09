package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/chayev/yurl/yurllib"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Listening on :8080...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/results", viewResults)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// PageOutput : The contents and URL parameters that are exported
type PageOutput struct {
	Content string
	URL     string
	Prefix  string
	Bundle  string
}

func handler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("tpl/base.html")
	t.Execute(w, nil)
}

func viewResults(w http.ResponseWriter, r *http.Request) {

	var url string
	var prefix string
	var bundle string

	for _, n := range r.URL.Query()["url"] {
		url = n
	}

	for _, n := range r.URL.Query()["prefix"] {
		prefix = n
	}

	for _, n := range r.URL.Query()["bundle"] {
		bundle = n
	}

	var output []string

	if url == "" {
		output = append(output, "Enter URL to validate.")
	} else {
		output = yurllib.CheckDomain(url, prefix, bundle, true)
	}

	content := &PageOutput{URL: url, Prefix: prefix, Bundle: bundle}

	for _, item := range output {
		content.Content += item
	}

	t, _ := template.ParseFiles("tpl/results.html")
	t.Execute(w, &content)
}
