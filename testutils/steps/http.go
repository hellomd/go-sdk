package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/hellomd/go-sdk/testutils"
)

// HTTPRequest -
func HTTPRequest(method string, server *httptest.Server, uri string, body *gherkin.DocString, contentType string, headers map[string]string) (*http.Response, error) {
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

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return http.DefaultClient.Do(req)
}

// IGetFrom -
func IGetFrom(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return HTTPRequest("GET", server, uri, body, "", nil)
}

func IHeadFrom(server *httptest.Server, uri string) (*http.Response, error) {
	return HTTPRequest("HEAD", server, uri, nil, "", nil)
}

// IPostTo -
func IPostTo(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return HTTPRequest("POST", server, uri, body, "", nil)
}

// IPostToAs -
func IPostToAs(server *httptest.Server, uri string, as string, body *gherkin.DocString) (*http.Response, error) {
	return HTTPRequest("POST", server, uri, body, as, nil)
}

// IPutOn -
func IPutOn(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return HTTPRequest("PUT", server, uri, body, "", nil)
}

// IDeleteFrom -
func IDeleteFrom(server *httptest.Server, uri string, body *gherkin.DocString) (*http.Response, error) {
	return HTTPRequest("DELETE", server, uri, body, "", nil)
}

// TheStatusCodeShouldBe -
func TheStatusCodeShouldBe(response *http.Response, statusText string) error {
	if response == nil {
		return fmt.Errorf("Expected a recorded response")
	}

	if http.StatusText(response.StatusCode) != statusText {
		return fmt.Errorf("Expected status %v, but got %v instead", statusText, http.StatusText(response.StatusCode))
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
