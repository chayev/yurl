package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/chayev/yurl/yurllib"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Listening on :8080...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/results", viewResults)
	http.HandleFunc("/android", android)
	http.HandleFunc("/android-results", viewResultsAndroid)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// PageOutput : The contents and URL parameters that are exported
type PageOutput struct {
	Content     string
	URL         string
	Prefix      string
	Bundle      string
	CurrentTime time.Time
}

func handler(w http.ResponseWriter, r *http.Request) {

	content := &PageOutput{CurrentTime: time.Now()}

	t, _ := template.ParseFiles("tpl/home.html", "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToAndroid.html")
	t.ExecuteTemplate(w, "home.html", &content)
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
		output = yurllib.CheckAASADomain(url, prefix, bundle, true)
	}

	content := &PageOutput{URL: url, Prefix: prefix, Bundle: bundle}

	for _, item := range output {
		content.Content += item
	}

	content.CurrentTime = time.Now()

	t, _ := template.ParseFiles("tpl/results.html", "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToAndroid.html")
	t.ExecuteTemplate(w, "results.html", &content)
}

func android(w http.ResponseWriter, r *http.Request) {

	content := &PageOutput{CurrentTime: time.Now()}

	t, _ := template.ParseFiles("tpl/android.html", "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToiOS.html")
	t.ExecuteTemplate(w, "android.html", &content)
}

func viewResultsAndroid(w http.ResponseWriter, r *http.Request) {

	var url string
	var package_name string
	var fingerprint string

	for _, n := range r.URL.Query()["url"] {
		url = n
	}

	for _, n := range r.URL.Query()["prefix"] {
		package_name = n
	}

	for _, n := range r.URL.Query()["bundle"] {
		fingerprint = n
	}

	var output []string

	if url == "" {
		output = append(output, "Enter URL to validate.")
	} else {
		output = yurllib.CheckAssetLinkDomain(url, package_name, fingerprint)
	}

	content := &PageOutput{URL: url, Prefix: package_name, Bundle: fingerprint}

	for _, item := range output {
		content.Content += item
	}

	content.CurrentTime = time.Now()

	t, _ := template.ParseFiles("tpl/android-results.html", "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToiOS.html")
	t.ExecuteTemplate(w, "android-results.html", &content)
}
