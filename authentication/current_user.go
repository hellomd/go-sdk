package authentication

// CurrentUser -
type CurrentUser struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Empty -
func (cu *CurrentUser) Empty() bool {
	return cu.ID == ""
}
