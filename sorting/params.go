package sorting

import (
	"fmt"
	"strings"
)

const (
	// SortQueryParam -
	SortQueryParam = "sort"
)

// Extract -
func Extract(query map[string][]string, sorter Sorter) error {
	validFields := sorter.GetValidFields()
	for _, v := range query[SortQueryParam] {
		if _, ok := validFields[strings.Trim(v, "-")]; !ok {
			return fmt.Errorf("Invalid parameter in sort query string field %v", v)
		}
	}
	sorter.SetFields(query[SortQueryParam])
	return nil
}
