package sorting

type basicSorter struct {
	validFields map[string]bool
	fields      []string
}

// NewBasicSorter -
func NewBasicSorter(validFields []string) Sorter {
	vfMap := make(map[string]bool)
	for _, v := range validFields {
		vfMap[v] = true
	}
	return &basicSorter{
		validFields: vfMap,
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
