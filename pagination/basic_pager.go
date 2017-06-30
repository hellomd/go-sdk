package pagination

type basicPager struct {
	page    int
	perPage int
	url     string
}

// NewBasicPager -
func NewBasicPager(url string, defaultPerPage int) Pager {
	return &basicPager{
		page:    1,
		url:     url,
		perPage: defaultPerPage,
	}
}

func (p *basicPager) SetPage(page int) {
	p.page = page
}

func (p *basicPager) SetPerPage(perPage int) {
	p.perPage = perPage
}

func (p *basicPager) GetPage() int {
	return p.page
}

func (p *basicPager) GetPerPage() int {
	return p.perPage
}

func (p *basicPager) GetNextPage() int {
	return p.page + 1
}

func (p *basicPager) GetURL() string {
	return p.url
}
