package pagination

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	// TotalCountHeaderKey -
	TotalCountHeaderKey = "X-Total-Count"
	// LinkHeaderKey -
	LinkHeaderKey = "Link"
	// LinkTemplate -
	LinkTemplate = "<%s?%s>; rel=\"next\""
)

// SetTotalHeader -
func SetTotalHeader(h http.Header, count int) {
	h.Set(TotalCountHeaderKey, strconv.Itoa(count))
}

// SetLinkHeader -
func SetLinkHeader(h http.Header, query map[string][]string, pager Pager) {
	x := []string{}

	query["page"] = []string{strconv.Itoa(pager.GetNextPage())}

	for k, v := range query {
		if len(v) == 1 {
			x = append(x, fmt.Sprintf("%s=%s", k, v[0]))
		} else {
			x = append(x, fmt.Sprintf("%s=[%s]", k, strings.Join(v, ", ")))
		}
	}

	sort.Strings(x)

	h.Set(LinkHeaderKey, fmt.Sprintf(LinkTemplate, pager.GetURL(), strings.Join(x, "&")))
}
