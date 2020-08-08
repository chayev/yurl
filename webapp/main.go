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
	http.HandleFunc("/main", viewMain)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type PageOutput struct {
	Content string
}

func handler(w http.ResponseWriter, r *http.Request) {

	input := "https://kohls.onelink.me/asdas"

	output := yurllib.CheckDomain(input, "", "", true)

	var content PageOutput

	for _, item := range output {
		content.Content += item
	}

	t, _ := template.ParseFiles("tpl/base.html")
	t.Execute(w, &content)
}

func viewMain(w http.ResponseWriter, r *http.Request) {

	input := "https://kohls.onelink.me/asdas"

	output := yurllib.CheckDomain(input, "", "", true)

	var content PageOutput

	for _, item := range output {
		content.Content += item
	}

	t, _ := template.ParseFiles("tpl/main.html")
	t.Execute(w, &content)
}
