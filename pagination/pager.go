package pagination

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// TotalCountHeaderKey -
	TotalCountHeaderKey = "X-Total-Count"
	// LinkHeaderKey -
	LinkHeaderKey = "Link"

	pageParam    = "page"
	perPageParam = "perPage"
)

// NewPager creates a new Pager by extracting query parameters from a HTTP request
func NewPager(r *http.Request, defaultPerPage, maxPerPage int) *Pager {
	query := r.URL.Query()

	page := 1
	if v, err := strconv.Atoi(query.Get(pageParam)); err == nil && v > 0 {
		page = v
	}

	perPage := defaultPerPage
	if v, err := strconv.Atoi(query.Get(perPageParam)); err == nil && v > 0 {
		if v > maxPerPage {
			perPage = maxPerPage
		} else {
			perPage = v
		}
	}

	return &Pager{r.URL, page, perPage}
}

// Pager -
type Pager struct {
	CurrentURL *url.URL
	Page       int
	PerPage    int
}

// SetHeaders sets response headers for the current pager's state
func (p *Pager) SetHeaders(h http.Header, totalCount int) {
	h.Set(TotalCountHeaderKey, strconv.Itoa(totalCount))
	h.Del(LinkHeaderKey)

	if p.Page > 1 {
		p.addLink(h, 1, "first")
		p.addLink(h, p.Page-1, "prev")
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(p.PerPage)))
	if p.Page < lastPage {
		p.addLink(h, p.Page+1, "next")
		p.addLink(h, lastPage, "last")
	}
}

func (p *Pager) addLink(h http.Header, page int, rel string) {
	uri := *p.CurrentURL
	q := uri.Query()
	q.Set(pageParam, strconv.Itoa(page))
	uri.RawQuery = q.Encode()
	h.Add(LinkHeaderKey, fmt.Sprintf(`<%v>; rel="%v"`, uri.String(), rel))
}
