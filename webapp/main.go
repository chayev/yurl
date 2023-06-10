package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chayev/yurl/yurllib"
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route requests to their corresponding handlers
	http.HandleFunc("/", homeHandler)                       // Home Page
	http.HandleFunc("/ios", formHandler)                    // iOS validation page
	http.HandleFunc("/android", formHandler)                // Android validation page
	http.HandleFunc("/ios-results", viewResultsHandler)     // Validation results page for iOS
	http.HandleFunc("/android-results", viewResultsHandler) // Validation results page for Android

	// Start the HTTP server
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var envRoot = os.Getenv("Y_THEME_ROOT")

// Initialize the templates on program start-up
var templates = template.Must(template.ParseFiles(
	envRoot + "tpl/home.html",
	envRoot + "tpl/ios.html",
	envRoot + "tpl/ios-results.html",
	envRoot + "tpl/android.html",
	envRoot + "tpl/android-results.html",
	envRoot + "tpl/partials/header.html",
	envRoot + "tpl/partials/footer.html",
	envRoot + "tpl/partials/navToAndroid.html",
	envRoot + "tpl/partials/navToiOS.html",
	envRoot + "tpl/partials/copyLink.html",
))

// PageOutput : The contents and URL parameters that are exported
type PageOutput struct {
	Content     string
	URL         string
	Prefix      string
	Bundle      string
	CurrentTime time.Time
}

// Handler function for the home pages
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize PageOutput with the current time
	content := &PageOutput{CurrentTime: time.Now()}

	// Render the template and handle errors
	err := templates.ExecuteTemplate(w, "home.html", content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

// Handler function for the iOS and Android validation pages
func formHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize PageOutput with the current time
	content := &PageOutput{CurrentTime: time.Now()}

	// Determine if the request is for Android validation
	isAndroid := r.URL.Path == "/android"

	var templateName string
	if isAndroid {
		templateName = "android.html"
	} else {
		templateName = "ios.html"
	}

	// Render the template and handle errors
	err := templates.ExecuteTemplate(w, templateName, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

// Handler function for the iOS and Android validation results pages
func viewResultsHandler(w http.ResponseWriter, r *http.Request) {
	var url string
	var prefix string
	var bundle string
	var isAndroid bool

	// Determine if the request is for iOS or Android validation
	if r.URL.Path == "/ios-results" {
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
	url = r.FormValue("url")
	prefix = r.FormValue("prefix")
	bundle = r.FormValue("bundle")

	var output []string

	// Check if the URL parameter is empty, if so add a message to the output
	if url == "" {
		output = append(output, "Enter URL to validate.")
	} else {
		log.Println("\n####Validating url: " + url + "\n")
		// Call the appropriate validation function based on the isAndroid boolean
		if isAndroid {
			output = yurllib.CheckAssetLinkDomain(url, prefix, bundle)
		} else {
			output = yurllib.CheckAASADomain(url, prefix, bundle, true)
		}
	}

	// Initialize the PageOutput struct with the URL, prefix, and bundle
	content := &PageOutput{URL: url, Prefix: prefix, Bundle: bundle}

	// Add each item in the output slice to the PageOutput struct's Content field
	for _, item := range output {
		content.Content += item
	}

	// Set the CurrentTime field of the PageOutput struct to the current time
	content.CurrentTime = time.Now()

	var templateName string

	// Determine which template to use based on the isAndroid boolean
	if isAndroid {
		templateName = "android-results.html"
	} else {
		templateName = "ios-results.html"
	}

	// Render the template and handle errors
	err = templates.ExecuteTemplate(w, templateName, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
