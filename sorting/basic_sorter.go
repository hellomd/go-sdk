package sorting

type basicSorter struct {
	validFields map[string]bool
	fields      []string
}

// NewBasicSorter -
func NewBasicSorter(validFields map[string]bool) Sorter {
	return &basicSorter{
		validFields: validFields,
	}
}

func (p *basicSorter) GetFields() []string {
	return p.fields
}

func (p *basicSorter) SetFields(fields []string) {
	p.fields = fields
}

func (p *basicSorter) GetValidFields() map[string]bool {
	return p.validFields
}
