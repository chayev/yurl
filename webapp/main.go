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
	http.HandleFunc("/android", handler)
	http.HandleFunc("/results", viewResults)
	http.HandleFunc("/android-results", viewResults)
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

	// Determine if the request is for Android validation
	isAndroid := r.URL.Path == "/android"

	var templateName string
	if isAndroid {
		templateName = "android.html"
	} else {
		templateName = "home.html"
	}

	t, _ := template.ParseFiles("tpl/"+templateName, "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToAndroid.html", "tpl/partials/navToiOS.html")
	t.ExecuteTemplate(w, templateName, &content)
}

func viewResults(w http.ResponseWriter, r *http.Request) {

	var url string
	var prefix string
	var bundle string
	var isAndroid bool

	// Determine if the request is for iOS or Android validation
	if r.URL.Path == "/results" {
		isAndroid = false
	} else if r.URL.Path == "/android-results" {
		isAndroid = true
	} else {
		// Invalid request path
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	// Parse form data from request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get URL, prefix, and bundle parameters from form data
	url = r.Form.Get("url")
	prefix = r.Form.Get("prefix")
	bundle = r.Form.Get("bundle")

	var output []string

	if url == "" {
		output = append(output, "Enter URL to validate.")
	} else {
		if isAndroid {
			output = yurllib.CheckAssetLinkDomain(url, prefix, bundle)
		} else {
			output = yurllib.CheckAASADomain(url, prefix, bundle, true)
		}
	}

	content := &PageOutput{URL: url, Prefix: prefix, Bundle: bundle}

	for _, item := range output {
		content.Content += item
	}

	content.CurrentTime = time.Now()

	var templateName string

	if isAndroid {
		templateName = "android-results.html"
	} else {
		templateName = "results.html"
	}

	t, _ := template.ParseFiles("tpl/"+templateName, "tpl/partials/header.html", "tpl/partials/footer.html", "tpl/partials/navToAndroid.html", "tpl/partials/navToiOS.html")

	t.ExecuteTemplate(w, templateName, &content)
}
