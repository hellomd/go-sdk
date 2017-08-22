package authentication

// User -
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsService   bool `json:"isService"`
}

// Empty -
func (u *User) Empty() bool {
	return u.ID == ""
}
