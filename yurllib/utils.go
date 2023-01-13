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
	}

	scheme := parsedURL.Scheme

	if scheme != "https" {
		output = append(output, fmt.Sprintf("WARNING: The URL must use HTTPS, trying HTTPS instead. \n\n"))

		parsedURL.Scheme = "https"
		parsedURL, err = url.Parse(parsedURL.String())
		if err != nil {
			output = append(output, fmt.Sprintf("The URL failed to parse with error %s \n", err))
		}
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
