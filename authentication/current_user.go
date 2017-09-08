package authentication

// User -
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsService *bool  `json:"isService,omitempty"`
}

// Empty -
func (u *User) Empty() bool {
	return u.ID == ""
}
