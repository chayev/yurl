package yurllib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type assetLinkFile []struct {
	Target   Target   `json:"target,omitempty"`
	Relation []string `json:"relation,omitempty"`
}
type Target struct {
	PackageName            string   `json:"package_name,omitempty"`
	Sha256CertFingerprints []string `json:"sha256_cert_fingerprints,omitempty"`
	Namespace              string   `json:"namespace,omitempty"`
}

//https://developer.android.com/training/app-links/verify-android-applinks#publish-json
//https://developers.google.com/digital-asset-links/reference/rest/v1/Asset
//https://github.com/google/digitalassetlinks/blob/master/well-known/details.md
//Above defines requirements for Asset Links

// CheckAssetLinkDomain : Main function used by CLI and WebApp for Android App Link validation
func CheckAssetLinkDomain(inputURL string, packageInput string, fingerprintInput string) []string {

	var output []string

	cleanedDomain, messages := getDomain(inputURL)

	output = append(output, messages...)

	rawResult, message, errors := loadAssetLinkContents(cleanedDomain)
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

	if contentType[0] != "application/json" {
		output = append(output, fmt.Sprint("\nInvalid content type. Expecting [application/json]. Please update and test again. \n"))
		output = append(output, fmt.Sprint("\nIf you believe this error is invalid, please open an issue on github or email support@chayev.com and we will investigate."))
		return output
	}

	result, err := ioutil.ReadAll(rawResult.Body)
	if err != nil {
		return output
	}

	messages, errorsx := evaluateAssetLink(result, packageInput, fingerprintInput)
	if len(errorsx) > 0 {
		output = append(output, fmt.Sprintf("\n  %s\n", "Errors: "))
		for _, e := range errorsx {
			output = append(output, fmt.Sprintf("  - %s\n", e))
		}

		output = append(output, fmt.Sprintf("\n  %s\n", "File Contents: "))
		output = append(output, fmt.Sprintf("\n  %s\n", result))
		return output
	}

	output = append(output, messages...)

	return output
}

func loadAssetLinkContents(domain string) (*http.Response, []string, []error) {

	var output []string
	var formatErrors []error
	var respStatus int

	wellKnownPath := "https://" + domain + "/.well-known/assetlinks.json"

	// Testing URLs
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-empty.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-invalidnamespace.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-nopackage.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-nofp.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-emptyfp.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-malformed.json"
	// wellKnownPath = "https://chayev.github.io/webpage-sandbox/assetlinks-onlynamespace.json"

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

	formatErrors = append(formatErrors, errors.New("It looks like this domain does not have Android Asset Links configured. No json file found in expected location."))

	return nil, output, formatErrors
}

func evaluateAssetLink(result []byte, packageInput string, fingerprintInput string) ([]string, []error) {

	var output []string
	var formatErrors []error

	var reqResp assetLinkFile

	err := json.Unmarshal(result, &reqResp)
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

	output = append(output, fmt.Sprintln("JSON Validation: \t\t  Pass"))

	validJSON, formatErrors := verifyAssetLinkJSONformat(reqResp)

	prettyJSON, err := json.MarshalIndent(reqResp, "", "    ")
	if err != nil {
		formatErrors = append(formatErrors, fmt.Errorf("ioutil.ReadAll failed to parse with error: \n%w", err)) //define this better
		return output, formatErrors
	}

	if validJSON {
		output = append(output, fmt.Sprintln("JSON Schema: \t\t\t  Pass"))

		if packageInput != "" && fingerprintInput != "" {
			if verifyInputIsPresent(reqResp, packageInput, fingerprintInput) {
				output = append(output, fmt.Sprintln("Package/Fingerprint availability: Pass"))
			} else {
				output = append(output, fmt.Sprintln("Package/Fingerprint availability: Fail"))
			}
		}

		output = append(output, fmt.Sprintf("\n%s\n", string(prettyJSON)))

	} else {
		output = append(output, fmt.Sprintln("\nJSON Schema: Fail"))
		output = append(output, fmt.Sprintln("\nTarget must have an 'android_app' namespace with both a 'package_name' and atleast one 'sha256_cert_fingerprints' in an array."))

		for _, formatError := range formatErrors {
			output = append(output, fmt.Sprintf("\n  %s\n", "Errors: "))
			output = append(output, fmt.Sprintf("  - %s\n", formatError))
		}

		output = append(output, fmt.Sprintf("\n%s\n", string(prettyJSON)))

		return output, formatErrors
	}

	return output, formatErrors

}

func verifyAssetLinkJSONformat(content assetLinkFile) (bool, []error) {

	assetLinks := content //Array

	var formatErrors []error

	isValid := false

	hasAppNamespace := false
	hasPackageName := false
	hasFingerprint := false

	if assetLinks == nil {
		formatErrors = append(formatErrors, errors.New("No data found in the file."))
	} else {
		for _, assetLink := range assetLinks {
			namespace := assetLink.Target.Namespace
			packageName := assetLink.Target.PackageName
			fingerprints := assetLink.Target.Sha256CertFingerprints

			if namespace == "android_app" {
				hasAppNamespace = true

				if len(packageName) != 0 {
					hasPackageName = true

					for _, fingerprint := range fingerprints {
						if len(fingerprint) != 0 {
							hasFingerprint = true
							isValid = true
						}
					}
				}
			}
			if len(packageName) != 0 {
				hasPackageName = true
			}

			for _, fingerprint := range fingerprints {
				if len(fingerprint) != 0 {
					hasFingerprint = true
				}
			}
		}
	}

	if !hasAppNamespace {
		formatErrors = append(formatErrors, errors.New("None of the targets contains a namespace for 'android_app'."))
	}

	if !hasPackageName {
		formatErrors = append(formatErrors, errors.New("Target with 'android_app' namespace must have 'package_name'."))
	}

	if !hasFingerprint {
		formatErrors = append(formatErrors, errors.New("Target with 'android_app' namespace must have 'sha256_cert_fingerprints'."))
	}

	return isValid, formatErrors

}

func verifyInputIsPresent(content assetLinkFile, packageInput string, fingerprintInput string) bool {

	for i := 0; i < len(content); i++ {
		record := content[i].Target
		if record.Namespace == "android_app" {
			if record.PackageName == packageInput {
				for j := 0; j < len(record.Sha256CertFingerprints); j++ {
					fingerprint := record.Sha256CertFingerprints[j]
					if fingerprint == fingerprintInput {
						return true
					}
				}
			}
		}
	}

	return false
}
