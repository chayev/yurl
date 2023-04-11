package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(homeHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := `<h1><span id="name">yURL</span>: Deep Linking Validator</h1>`
    if !strings.Contains(rr.Body.String(), expected) {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestFormHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(formHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := `<h1><span id="name"><a href="/" aria-label="Go to Homepage" tabindex="-1" style="color:black">yURL</a></span>: Universal Links / AASA File Validator</h1>`
    if !strings.Contains(rr.Body.String(), expectedBody) {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedBody)
    }
}

func TestFormHandlerAndroid(t *testing.T) {
	req, err := http.NewRequest("GET", "/android", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(formHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := `<h1><span id="name"><a href="/" aria-label="Go to Homepage" tabindex="-1" style="color:black">yURL</a></span>: Asset Links File Validator</h1>`
    if !strings.Contains(rr.Body.String(), expectedBody) {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedBody)
    }
}

func TestViewResultsHandler(t *testing.T) {
    data := url.Values{}
    data.Set("url", "https://suadeo.onelink.me")
    data.Set("prefix", "com.example")
    data.Set("bundle", "com.example.app")
    body := strings.NewReader(data.Encode())

    req, err := http.NewRequest("POST", "/ios-results", body)
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(viewResultsHandler)

    handler.ServeHTTP(rr, req)
	
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedBody := `Found file at:
  https://suadeo.onelink.me/.well-known/apple-app-site-association`
    if !strings.Contains(rr.Body.String(), expectedBody) {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedBody)
    }
}

func TestViewResultsHandlerAndroid(t *testing.T) {
    data := url.Values{}
    data.Set("url", "https://suadeo.onelink.me")
    data.Set("prefix", "com.example")
    data.Set("bundle", "com.example.app")
    body := strings.NewReader(data.Encode())

    req, err := http.NewRequest("POST", "/android-results", body)
    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(viewResultsHandler)

    handler.ServeHTTP(rr, req)
	
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedBody := `Found file at:
  https://suadeo.onelink.me/.well-known/assetlinks.json`
    if !strings.Contains(rr.Body.String(), expectedBody) {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expectedBody)
    }
}
