package model

type User struct {
	Username       string              `json:"username"`
	Password       string              `json:"password"`
	Authentication *smtpAuthentication `json:"authentication,omitempty"`
}
