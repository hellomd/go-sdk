package pagination

// Pager -
type Pager interface {
	SetPage(int)
	SetPerPage(int)
	GetPage() int
	GetPerPage() int
	GetNextPage() int
	GetURL() string
}
