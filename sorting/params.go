package sorting

import "fmt"

const (
	// SortQueryParam -
	SortQueryParam = "sort"
)

// Collect -
func Collect(query map[string][]string, sorter Sorter) error {
	validFields := sorter.GetValidFields()
	for _, v := range query[SortQueryParam] {
		if _, ok := validFields[v]; !ok {
			return fmt.Errorf("Invalid parameter in sort query string field %v", v)
		}
	}
	sorter.SetFields(query[SortQueryParam])
	return nil
}
