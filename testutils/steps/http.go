package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/hellomd/middlewares/testutils"
)

func httpRequest(method string, server *httptest.Server, uri string, body *gherkin.DocString, contentType string) (*http.Response, error) {
	var content io.Reader
	if body != nil {
		content = strings.NewReader(body.Content)
	}

	req, err := http.NewRequest(
		method, server.URL+uri, content,
	)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}

	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Add("Content-Type", contentType)

	return http.DefaultClient.Do(req)
}

// IGetFrom -
func IGetFrom(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return httpRequest("GET", server, uri, body, "")
}

// IPostTo -
func IPostTo(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return httpRequest("POST", server, uri, body, "")
}

// IPostToAs -
func IPostToAs(server *httptest.Server, uri string, as string, body *gherkin.DocString) (*http.Response, error) {
	return httpRequest("POST", server, uri, body, as)
}

// IPutOn -
func IPutOn(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return httpRequest("PUT", server, uri, body, "")
}

// IDeleteFrom -
func IDeleteFrom(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return httpRequest("DELETE", server, uri, body, "")
}

// TheStatusCodeShouldBe -
func TheStatusCodeShouldBe(response *http.Response, statusText string) error {
	if response == nil {
		return fmt.Errorf("Expected a recorded response")
	}

	switch statusText {
	case "OK":
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("Expected status code 200 (OK), but got %v instead",
				response.StatusCode)
		}

	case "Not Found":
		if response.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Expected status code 404 (Not Found), but got %v instead",
				response.StatusCode)
		}

	case "Unprocessable Entity":
		if response.StatusCode != http.StatusUnprocessableEntity {
			return fmt.Errorf("Expected status code 422 (Unprocessable Entity), but got %v instead",
				response.StatusCode)
		}

	case "Bad Request":
		if response.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("Expected status code 400 (Bad Request), but got %v instead",
				response.StatusCode)
		}

	default:
		return fmt.Errorf("Invalid status code %s", statusText)
	}

	return nil
}

// TheJSONResponseShouldBe -
func TheJSONResponseShouldBe(response *http.Response, expectedJSON *gherkin.DocString) error {

	if response == nil {
		return fmt.Errorf("Expected a recorded response")
	}

	var expectedObj interface{}
	if err := json.Unmarshal([]byte(expectedJSON.Content), &expectedObj); err != nil {
		return err
	}

	var actualObj interface{}

	if err := json.NewDecoder(response.Body).Decode(&actualObj); err != nil {
		return err
	}

	if !testutils.JSONEqualsIgnoreOrder(actualObj, expectedObj) {
		actualJSON, _ := json.MarshalIndent(actualObj, "", "  ")
		expJSON, _ := json.MarshalIndent(expectedObj, "", "  ")
		return fmt.Errorf("Expected:\n%s\n\nGot JSON\n%s", string(expJSON), string(actualJSON))
	}

	return nil
}
