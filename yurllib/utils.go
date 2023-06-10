package yurllib

import (
	"fmt"
	"net/http"
	"net/url"
)

func getDomain(input string) (string, []string) {

	var output []string

	//Clean up domains, removing scheme and path
	parsedURL, err := url.Parse(input)
	if err != nil {
		output = append(output, fmt.Sprintf("The URL failed to parse with error %s \n", err))
		return "", output
	}

	scheme := parsedURL.Scheme

	if scheme != "https" && scheme != "" {
		output = append(output, fmt.Sprintf("WARNING: The URL must use HTTPS, changing the protocol to HTTPS instead. \n\n"))
	}
	
	parsedURL.Scheme = "https"
	parsedURL, err = url.Parse(parsedURL.String())
	if err != nil {
		output = append(output, fmt.Sprintf("The URL failed to parse with error %s \n", err))
		return "", output
	}

	return parsedURL.Host, output
}

func makeRequest(fileURL string) (*http.Response, error) {

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CompareStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func OutputStatus(output []string, name string, pass bool) []string{
	if pass {		
		return append(output, fmt.Sprintf("%s: Pass\n", name))
	} else {
		return append(output, fmt.Sprintf("%s: Fail\n", name))
	}
}
