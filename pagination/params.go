package pagination

import (
	"errors"
	"strconv"
)

const (
	// PageQueryParam -
	PageQueryParam = "page"
	// PerPageQueryParam -
	PerPageQueryParam = "perPage"
)

// ErrMaxPerPageExceeded -
var ErrMaxPerPageExceeded = errors.New("Max per page exceeded")

func extractParam(key string, query map[string][]string) (int, error) {
	if len(query[key]) > 0 {
		rs, err := strconv.Atoi(query[key][0])
		if err != nil {
			return 0, err
		}
		return rs, nil
	}
	return 0, nil
}

// ExtractPage -
func ExtractPage(query map[string][]string, pager Pager) error {
	page, err := extractParam(PageQueryParam, query)
	if err != nil {
		return err
	}
	if page != 0 {
		pager.SetPage(page)
	}
	return nil
}

// ExtractPerPage -
func ExtractPerPage(query map[string][]string, pager Pager) error {
	maxPerPage := pager.GetMaxPerPage()
	perPage, err := extractParam(PerPageQueryParam, query)
	if err != nil {
		return err
	}
	if perPage > maxPerPage {
		return ErrMaxPerPageExceeded
	}
	if perPage != 0 {
		pager.SetPerPage(perPage)
	}
	return nil
}

// Extract -
func Extract(query map[string][]string, pager Pager) error {
	err := ExtractPage(query, pager)
	if err != nil {
		return err
	}

	err = ExtractPerPage(query, pager)
	if err != nil {
		return err
	}
	return nil
}
