package pagination

// Pager -
type Pager interface {
	SetPage(int)
	SetPerPage(int)
	GetNextPage() int
	GetPerPage() int
	GetURL() string
}
