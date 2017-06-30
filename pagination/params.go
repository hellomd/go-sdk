package pagination

import (
	"strconv"
)

const (
	// PageQueryParam -
	PageQueryParam = "page"
	// PerPageQueryParam -
	PerPageQueryParam = "perPage"
)

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

// GetPage -
func GetPage(query map[string][]string, pager Pager) error {
	page, err := collectParam(PageQueryParam, query)
	if err != nil {
		return err
	}
	if page != 0 {
		pager.SetPage(page)
	}
	pager.SetPage(page)
	return nil
}

// GetPerPage -
func GetPerPage(query map[string][]string, pager Pager) error {
	perPage, err := collectParam(PerPageQueryParam, query)
	if err != nil {
		return err
	}
	if perPage != 0 {
		pager.SetPerPage(perPage)
	}
	return nil
}
