package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/chayev/yurl/yurllib"
)

func main() {
	http.HandleFunc("/", handler)
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
