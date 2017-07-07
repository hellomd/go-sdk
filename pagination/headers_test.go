package pagination

import (
	"net/http"
	"testing"
)

// SetLinkHeader -
// func SetLinkHeader(h http.Header, pager Pager) {
// 	h.Set(LinkHeaderKey, fmt.Sprintf(LinkTemplate, pager.GetURL(), pager.GetNextPage(), pager.GetPerPage()))
// }

const url = "hellomd.io"

func getPager() Pager {
	return NewBasicPager(url, 20, 50)
}

func TestSetLinkHeader(t *testing.T) {
	pager := getPager()
	expectedURL := "<" + url + "?page=2>; rel=\"next\""

	header := http.Header{}
	SetLinkHeader(header, map[string][]string{}, pager)

	actualURL := header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}
}

func TestSetLinkHeaderWithPerPage(t *testing.T) {
	pager := getPager()
	pager.SetPage(2)
	pager.SetPerPage(10)
	expectedURL := "<" + url + "?page=3&perPage=10>; rel=\"next\""

	header := http.Header{}
	SetLinkHeader(header, map[string][]string{"page": []string{"2"}, "perPage": []string{"10"}}, pager)

	actualURL := header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}
}

func TestSetLinkHeaderNextPage(t *testing.T) {
	pager := getPager()
	pager.SetPage(2)
	expectedURL := "<" + url + "?page=3>; rel=\"next\""

	header := http.Header{}
	SetLinkHeader(header, map[string][]string{}, pager)

	actualURL := header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}

	SetLinkHeader(header, map[string][]string{"page": []string{"3"}}, pager)

	actualURL = header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}
}

func TestSetLinkHeaderPreservesOtherParams(t *testing.T) {
	pager := getPager()
	pager.SetPage(2)
	expectedURL := "<" + url + "?page=3&something=else>; rel=\"next\""

	header := http.Header{}
	SetLinkHeader(header, map[string][]string{"page": []string{"3"}, "something": []string{"else"}}, pager)

	actualURL := header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}
}

func TestSetLinkHeaderPreservesOtherArrayParams(t *testing.T) {
	pager := getPager()
	pager.SetPage(2)
	expectedURL := "<" + url + "?page=3&something=[else, more]>; rel=\"next\""

	header := http.Header{}
	SetLinkHeader(header, map[string][]string{"page": []string{"3"}, "something": []string{"else", "more"}}, pager)

	actualURL := header.Get(LinkHeaderKey)
	if actualURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, actualURL)
	}
}
