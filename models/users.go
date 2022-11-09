package models

type User struct {
	Username       string              `json:"username"`
	Password       string              `json:"password"`
	Authentication *smtpAuthentication `json:"authentication,omitempty"`
}

// InvalidPassword returns true if the given password does not match.
func (m *User) InvalidPassword(password string) bool {

	if password == "" {
		return true
	}

	if m.Password != password {
		return true
	}

	return false
}
