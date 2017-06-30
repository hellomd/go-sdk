package sorting

// Sorter -
type Sorter interface {
	GetValidFields() map[string]bool
	GetFields() []string
	SetFields([]string)
}
