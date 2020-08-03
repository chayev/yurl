package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <URL>",
	Short: "Validate your link against Apple's requirements",
	Run: func(cmd *cobra.Command, args []string) {
		checkDomain(args[0], "", "", true)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type detail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}
type appLinks struct {
	Apps    []string `json:"apps"`
	Details []detail `json:"details"`
}

type aasaFile struct {
	Applinks *appLinks `json:"applinks"`
}

func checkDomain(inputURL string, bundleIdentifier string, teamIdentifier string, allowUnencrypted bool) {
	//Clean up domains, removing scheme and path
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		fmt.Printf("The URL failed to parse with error %s \n", err)
	}

	cleanedDomain := parsedURL.Host
	scheme := parsedURL.Scheme

	if scheme != "https" {
		fmt.Printf("WARNING: The URL must use HTTPS, trying HTTPS instead. \n\n")
	}

	// fmt.Println(cleanedDomain)

	// call loadAASAContents and handle response
	result, err := loadAASAContents(cleanedDomain)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	defer result.Body.Close()

	contentType := result.Header["Content-Type"]

	isEncryptedMimeType := contentType[0] == "application/pkcs7-mime"
	isJSONMimeType := contentType[0] == "application/json" || contentType[0] == "text/json"
	isJSONTypeOK := allowUnencrypted && isJSONMimeType // Only ok if both the "allow" flag is true, and... it's a valid type.

	if !isEncryptedMimeType && !isJSONTypeOK {
		fmt.Printf("Invalid content-type: %s \n", contentType[0])
		//return nil or error
	}

	if allowUnencrypted {
		// Try to decode the JSON right away (this assumes the file is not encrypted)
		// If it's not encrypted, we'll just return it
		evaluateAASA(result, bundleIdentifier, teamIdentifier, false)

	} else {
		// Decrypt and evaluate file
	}

}

func loadAASAContents(domain string) (*http.Response, error) {
	wellKnownPath := "https://" + domain + "/.well-known/apple-app-site-association"
	aasaPath := "https://" + domain + "/apple-app-site-association"

	resp := makeRequest(wellKnownPath)
	respStatus := resp.StatusCode

	if respStatus >= 200 && respStatus < 300 {
		fmt.Printf("Found file at:\n  %s\n\n", wellKnownPath)
		fmt.Println("No Redirect: Pass")
		return resp, nil
	}

	resp = makeRequest(aasaPath)
	respStatus = resp.StatusCode

	if respStatus >= 200 && respStatus < 300 {
		fmt.Printf("Found file at:\n  %s\n\n", aasaPath)
		fmt.Println("No Redirect: Pass")
		return resp, nil
	}

	return nil, errors.New("could not find file in either known locations")
}

func makeRequest(fileURL string) *http.Response {
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Fatal("The http request failed with error", err)
		//Handle Error messaging better
	}
	// defer resp.Body.Close()

	return resp
}

func evaluateAASA(result *http.Response, bundleIdentifier string, teamIdentifier string, encrypted bool) {

	jsonText, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll failed to parse with error", err) //define this better
	}

	var reqResp aasaFile

	err = json.Unmarshal(jsonText, &reqResp)
	if err != nil {
		prettyJSON, err := json.MarshalIndent(jsonText, "", "    ")
		if err != nil {
			log.Fatal("Failed to print contents", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))

		log.Fatal("JSON Validation: Fail")
	}

	fmt.Println("JSON Validation: Pass")

	validJSON, formatErrors := verifyJSONformat(reqResp)

	if validJSON {
		fmt.Printf("JSON Schema: Pass\n\n")

		prettyJSON, err := json.MarshalIndent(reqResp, "", "    ")
		if err != nil {
			log.Fatal("Failed to print contents", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))

	} else {
		fmt.Println("JSON Schema: Fail")
		for _, formatError := range formatErrors {
			fmt.Printf("  %s\n", formatError)
		}
	}

	if validJSON && bundleIdentifier != "" {
		if verifyBundleIdentifierIsPresent(reqResp, bundleIdentifier, teamIdentifier) {
			fmt.Println("Team/Bundle availability: Pass")
		} else {
			fmt.Println("Team/Bundle availability: Fail")
		}

	}

}

func verifyJSONformat(content aasaFile) (bool, []error) {
	appLinks := content.Applinks

	var formatErrors []error

	if appLinks == nil {
		formatErrors = append(formatErrors, errors.New("missing applinks region"))
	}

	apps := appLinks.Apps
	if apps == nil {
		formatErrors = append(formatErrors, errors.New("missing applinks/apps region"))
	} else if len(apps) > 0 {
		formatErrors = append(formatErrors, errors.New("the apps key must have its value be an empty array"))
	}

	details := appLinks.Details
	if details == nil {
		formatErrors = append(formatErrors, errors.New("missing applinks/details region"))
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
	}

	return false
}
