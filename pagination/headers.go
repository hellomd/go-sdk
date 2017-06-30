package pagination

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	// TotalCountHeaderKey -
	TotalCountHeaderKey = "X-Total-Count"
	// LinkHeaderKey -
	LinkHeaderKey = "Link"
	// LinkTemplate -
	LinkTemplate = "%s?page=%d&perPage=%d"
)

// SetTotalHeader -
func SetTotalHeader(h *http.Header, count int) {
	h.Set(TotalCountHeaderKey, strconv.Itoa(count))
}

// SetLinkHeader -
func SetLinkHeader(h *http.Header, pager Pager) {
	h.Set(LinkHeaderKey, fmt.Sprintf(LinkTemplate, pager.GetURL(), pager.GetNextPage(), pager.GetPerPage()))
}
