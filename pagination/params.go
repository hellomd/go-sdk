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

func collectParam(key string, query map[string][]string) (int, error) {
	if len(query[key]) > 0 {
		rs, err := strconv.Atoi(query[key][0])
		if err != nil {
			return 0, err
		}
		return rs, nil
	}
	return 0, nil
}

// CollectPage -
func CollectPage(query map[string][]string, pager Pager) error {
	page, err := collectParam(PageQueryParam, query)
	if err != nil {
		return err
	}
	if page != 0 {
		pager.SetPage(page)
	}
	return nil
}

// CollectPerPage -
func CollectPerPage(query map[string][]string, pager Pager) error {
	maxPerPage := pager.GetMaxPerPage()
	perPage, err := collectParam(PerPageQueryParam, query)
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

// Collect -
func Collect(query map[string][]string, pager Pager) error {
	err := CollectPage(query, pager)
	if err != nil {
		return err
	}

	err = CollectPerPage(query, pager)
	if err != nil {
		return err
	}
	return nil
}
