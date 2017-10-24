package pagination

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	defaultPerPage = 10
	maxPerPage     = 50
)

func TestPager(t *testing.T) {
	Convey(fmt.Sprintf("create pager with default %v and max %v per page", defaultPerPage, maxPerPage), t, func() {
		type testCase struct {
			rawQuery                      string
			expectedPage, expectedPerPage int
			resultDescription             string
		}

		cases := []testCase{
			{"?page=5", 5, defaultPerPage, "use requested page"},
			{"?perPage=5", 1, 5, "use requested per page"},
			{"?perPage=100", 1, maxPerPage, "use max per page"},
			{"", 1, defaultPerPage, "use defaults"},
			{"?page=1", 1, defaultPerPage, "use defaults"},
			{"?page=0", 1, defaultPerPage, "use defaults"},
			{"?page=-1", 1, defaultPerPage, "use defaults"},
			{"?page=abc", 1, defaultPerPage, "use defaults"},
			{"?perPage=-1", 1, defaultPerPage, "use defaults"},
			{"?perPage=0", 1, defaultPerPage, "use defaults"},
			{"?perPage=abc", 1, defaultPerPage, "use defaults"},
		}

		for _, c := range cases {
			title := c.rawQuery
			if title == "" {
				title = "<empty>"
			}

			Convey(title+" should "+c.resultDescription, func() {
				uri, _ := url.Parse("http://my.site/some/resource" + c.rawQuery)
				r := &http.Request{URL: uri}

				pager := NewPager(r, defaultPerPage, maxPerPage)

				So(pager.CurrentURL, ShouldEqual, uri)
				So(pager.Page, ShouldEqual, c.expectedPage)
				So(pager.PerPage, ShouldEqual, c.expectedPerPage)
			})
		}
	})

	Convey("set pager headers", t, func() {
		Convey("on first page", func() {
			uri, _ := url.Parse("http://my.site/some/resource")
			pager := &Pager{CurrentURL: uri, Page: 1, PerPage: 25}

			h := http.Header{}

			pager.SetHeaders(h, 80)

			Convey("only next and last links should be set", func() {
				So(h[TotalCountHeaderKey], ShouldResemble, []string{"80"})
				So(h[LinkHeaderKey], ShouldResemble, []string{
					`<http://my.site/some/resource?page=2>; rel="next"`,
					`<http://my.site/some/resource?page=4>; rel="last"`,
				})
			})
		})

		Convey("on last page", func() {
			uri, _ := url.Parse("http://my.site/some/resource")
			pager := &Pager{CurrentURL: uri, Page: 4, PerPage: 25}

			h := http.Header{}

			pager.SetHeaders(h, 80)

			Convey("only first and prev links should be set", func() {
				So(h[TotalCountHeaderKey], ShouldResemble, []string{"80"})
				So(h[LinkHeaderKey], ShouldResemble, []string{
					`<http://my.site/some/resource?page=1>; rel="first"`,
					`<http://my.site/some/resource?page=3>; rel="prev"`,
				})
			})
		})

		Convey("on middle page", func() {
			uri, _ := url.Parse("http://my.site/some/resource")
			pager := &Pager{CurrentURL: uri, Page: 3, PerPage: 25}

			h := http.Header{}

			pager.SetHeaders(h, 80)

			Convey("all four links should be set", func() {
				So(h[TotalCountHeaderKey], ShouldResemble, []string{"80"})
				So(h[LinkHeaderKey], ShouldResemble, []string{
					`<http://my.site/some/resource?page=1>; rel="first"`,
					`<http://my.site/some/resource?page=2>; rel="prev"`,
					`<http://my.site/some/resource?page=4>; rel="next"`,
					`<http://my.site/some/resource?page=4>; rel="last"`,
				})
			})
		})
	})
}
