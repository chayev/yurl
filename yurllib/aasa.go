package yurllib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.mozilla.org/pkcs7"
)

type component map[string]interface{} // https://developer.apple.com/documentation/bundleresources/applinks/details/components

type substitutionVariable map[string][]string // https://developer.apple.com/documentation/bundleresources/applinks/substitutionvariables

// https://developer.apple.com/documentation/bundleresources/applinks/defaults
//
//	type defaultStruct struct {
//		CaseSensitive  bool `json:"caseSensitive,omitempty"`
//		PercentEncoded bool `json:"percentEncoded,omitempty"`
//	}
type defaultStruct map[string]interface{}

type detail struct {
	AppID      string        `json:"appID,omitempty"`
	Paths      []string      `json:"paths,omitempty"`
	AppIDs     []string      `json:"appIDs,omitempty"`
	Components []component   `json:"components,omitempty"`
	Defaults   defaultStruct `json:"defaults,omitempty"`
}

type appLinks struct {
	Apps                  []string               `json:"apps,omitempty"`
	Details               []detail               `json:"details"`
	SubstitutionVariables []substitutionVariable `json:"substitutionVariables,omitempty"`
	Defaults              defaultStruct          `json:"defaults,omitempty"`
}

type aasaFile struct {
	Applinks             *appLinks       `json:"applinks"`
	WebCredentials       *genericService `json:"webcredentials,omitempty"`
	ActivityContinuation *genericService `json:"activitycontinuation,omitempty"`
	AppClips             *genericService `json:"appclips,omitempty"`
}
type genericService struct {
	Apps []string `json:"apps,omitempty"`
}

// CheckAASADomain : Main function used by CLI and WebApp for Apple App Site Association validation
func CheckAASADomain(inputURL string, bundleIdentifier string, teamIdentifier string, allowUnencrypted bool) []string {

	var output []string

	cleanedDomain, messages := getDomain(inputURL)

	output = append(output, messages...)

	// call loadAASAContents and handle response
	rawResult, message, errors := loadAASAContents(cleanedDomain)
	if len(errors) > 0 {
		for _, e := range errors {
			output = append(output, fmt.Sprintf("  %s\n", e))
		}
		return output
	}
	defer rawResult.Body.Close()

	output = append(output, message...)

	contentType := rawResult.Header["Content-Type"]

	output = append(output, fmt.Sprintf("Content-type: \t\t\t  %s \n", contentType))

	isEncryptedMimeType := false
	isJSONTypeOK := false

	if len(contentType) > 0 {
		isEncryptedMimeType = contentType[0] == "application/pkcs7-mime"
		isJSONMimeType := contentType[0] == "application/json" || contentType[0] == "text/json" || contentType[0] == "text/plain" || strings.Contains(contentType[0], "application/json") || contentType[0] == "application/octet-stream" || contentType[0] == "binary/octet-stream"
		isJSONTypeOK = allowUnencrypted && isJSONMimeType // Only ok if both the "allow" flag is true, and... it's a valid type.
	} else {
		isJSONTypeOK = true
	}

	if !isEncryptedMimeType && !isJSONTypeOK {
		output = append(output, fmt.Sprintf("Invalid content-type: \t\t  %s \n", contentType))
		output = append(output, fmt.Sprint("\nIf you believe this error is invalid, please open an issue on github or email support@chayev.com and we will investigate."))
		return output
	}

	result, err := ioutil.ReadAll(rawResult.Body)
	if err != nil {
		// formatErrors = append(formatErrors, fmt.Errorf("ioutil.ReadAll failed to parse with error: \n%w", err)) //define this better
		// return output, formatErrors
		return output
	}

	if allowUnencrypted {
		// Try to decode the JSON right away (this assumes the file is not encrypted)
		// If it's not encrypted, we'll just return it
		messages, errors := evaluateAASA(result, contentType, bundleIdentifier, teamIdentifier, false)
		if len(errors) > 0 {
			for _, e := range errors {
				output = append(output, fmt.Sprintf("  %s\n", e))
			}
			return output
		}

		output = append(output, messages...)

	} else {
		// Decrypt and evaluate file
	}

	return output
}

func loadAASAContents(domain string) (*http.Response, []string, []error) {

	var output []string
	var formatErrors []error
	var respStatus int

	wellKnownPath := "https://" + domain + "/.well-known/apple-app-site-association"
	aasaPath := "https://" + domain + "/apple-app-site-association"

	resp, err := makeRequest(wellKnownPath)
	if err == nil {
		respStatus = resp.StatusCode

		if respStatus >= 200 && respStatus < 300 {
			output = append(output, fmt.Sprintf("Found file at:\n  %s\n\n", wellKnownPath))
			output = append(output, fmt.Sprintln("No Redirect: \t\t\t  Pass"))
			return resp, output, nil
		}
	} else {
		formatErrors = append(formatErrors, fmt.Errorf("Error: %w", err))
	}

	resp, err = makeRequest(aasaPath)
	if err == nil {
		respStatus = resp.StatusCode

		if respStatus >= 200 && respStatus < 300 {
			output = append(output, fmt.Sprintf("Found file at:\n  %s\n\n", aasaPath))
			output = append(output, fmt.Sprintln("No Redirect: \t\t\t Pass"))
			return resp, output, nil
		}
	} else {
		formatErrors = append(formatErrors, fmt.Errorf("Error: %w", err))
	}

	formatErrors = append(formatErrors, errors.New("It looks like this domain does not have Universal Links configured. No association file found."))

	return nil, output, formatErrors
}

func evaluateAASA(result []byte, contentType []string, bundleIdentifier string, teamIdentifier string, encrypted bool) ([]string, []error) {

	var output []string
	var formatErrors []error

	var reqResp aasaFile

	err := json.Unmarshal(result, &reqResp)
	if err != nil {

		if len(contentType) > 0 && contentType[0] == "application/pkcs7-mime" {
			jsonTextb, err := pkcs7.Parse(result)
			if err != nil {
				formatErrors = append(formatErrors, fmt.Errorf("PKCS7 Parse Fail: \n%w", err)) //define this better
				return output, formatErrors
			}

			jsonText := jsonTextb.Content

			err = json.Unmarshal(jsonText, &reqResp)
		} else {
			if err != nil {
				prettyJSON, err := json.MarshalIndent(result, "", "    ")
				if err != nil {
					formatErrors = append(formatErrors, fmt.Errorf("ioutil.ReadAll failed to parse with error: \n%w", err)) //define this better
					return output, formatErrors
				}
				output = append(output, fmt.Sprintln("JSON Validation: Fail"))

				output = append(output, fmt.Sprintf("%s\n", string(prettyJSON)))

				return output, formatErrors
			}
		}
	}

	output = append(output, fmt.Sprintln("JSON Validation: \t\t  Pass"))

	validJSON, formatErrors := verifyAASAJSONformat(reqResp)

	if validJSON {
		output = append(output, fmt.Sprintln("JSON Schema: \t\t\t  Pass"))

		if bundleIdentifier != "" {
			if verifyBundleIdentifierIsPresent(reqResp, bundleIdentifier, teamIdentifier) {
				output = append(output, fmt.Sprintln("Team/Bundle availability: Pass"))
			} else {
				output = append(output, fmt.Sprintln("Team/Bundle availability: Fail"))
			}
		}

		prettyJSON, err := json.MarshalIndent(reqResp, "", "    ")
		if err != nil {
			formatErrors = append(formatErrors, fmt.Errorf("ioutil.ReadAll failed to parse with error: \n%w", err)) //define this better
			return output, formatErrors
		}
		output = append(output, fmt.Sprintf("\n%s\n", string(prettyJSON)))

	} else {
		output = append(output, fmt.Sprintln("JSON Schema: Fail"))
		for _, formatError := range formatErrors {
			output = append(output, fmt.Sprintf("  %s\n", formatError))
		}
		return output, formatErrors
	}

	return output, formatErrors

}

func verifyAASAJSONformat(content aasaFile) (bool, []error) {

	appLinks := content.Applinks

	var formatErrors []error

	if appLinks == nil {
		formatErrors = append(formatErrors, errors.New("missing applinks region"))
	} else {
		apps := appLinks.Apps
		if len(apps) > 0 {
			formatErrors = append(formatErrors, errors.New("the apps key must have its value be an empty array"))
		}

		details := appLinks.Details
		if details == nil {
			formatErrors = append(formatErrors, errors.New("missing applinks/details region"))
		}
	}

	if len(formatErrors) > 0 {
		return false, formatErrors
	}

	return true, formatErrors

}

func verifyBundleIdentifierIsPresent(content aasaFile, bundleIdentifier string, teamIdentifier string) bool {

	details := content.Applinks.Details
	matcher := bundleIdentifier + "." + teamIdentifier

	for i := 0; i < len(details); i++ {
		var detail = details[i]
		if detail.AppID == matcher && len(detail.Paths) > 0 {
			return true
		}

		var arrAppids = detail.AppIDs
		for x := 0; x < len(arrAppids); x++ {
			if arrAppids[x] == matcher {
				return true
			}
		}
	}

	return false
}
