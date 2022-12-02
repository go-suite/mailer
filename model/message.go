package model

type message struct {
	From      string `json:"from,omitempty"`
	To        string `json:"to" example:"me@gmail.com,you@gmail.com"`
	Subject   string `json:"subject"`
	Body      string `json:"body,omitempty"`
	HtmlBody  string `json:"html_body,omitempty" example:"Hello!"`
	PlainBody string `json:"plain_body,omitempty" example:"<p>Hello!</p>"`
}
