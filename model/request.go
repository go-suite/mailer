package model

type smtpAuthentication struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Request struct {
	Message        message             `json:"message"`
	Authentication *smtpAuthentication `json:"authentication,omitempty"`
}
