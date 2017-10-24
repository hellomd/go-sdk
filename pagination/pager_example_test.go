package pagination_test

import (
	"encoding/json"
	"net/http"

	"github.com/hellomd/go-sdk/pagination"
)

const (
	defaultPerPage = 25
	maxPerPage     = 100
)

func PagerExample() {
	http.HandleFunc("/things", func(w http.ResponseWriter, r *http.Request) {
		pager := pagination.NewPager(r, defaultPerPage, maxPerPage)

		things, total := ListThings(pager.Page, pager.PerPage)

		pager.SetHeaders(w.Header(), total)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(things)
	})

	http.ListenAndServe(":3000", nil)
}

func ListThings(page, perPage int) ([]string, int) {
	return []string{"foo", "bar"}, 10
}
