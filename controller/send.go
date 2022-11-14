package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-suite/mailer/config"
	status "github.com/go-suite/mailer/http/status"
	"github.com/go-suite/mailer/model"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
)

func Send(c *gin.Context) {
	var r model.Request
	var err error

	// Retrieve data sent in the body
	err = c.ShouldBind(&r)
	if err != nil {
		status.Error(c, http.StatusBadRequest, err)
		return
	}

	if config.C.Secure {
		// Retrieve token metadata
		td, err := tokenManager.ExtractTokenMetadata(c.Request)
		if err != nil {
			status.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Get smtp authentication from config if not specified in the body
		if r.Authentication == nil {
			u := config.C.GetUser(td.UserName)
			if u != nil {
				r.Authentication = u.Authentication
			}
		}
		if r.Authentication == nil {
			status.Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	err = sendMail(r)
	if err != nil {
		status.Error(c, http.StatusInternalServerError, err)
		return
	}

	status.Message(c, http.StatusOK, "Ok")
}

func sendMail(r model.Request) error {
	// decompose list of recipients
	tos := strings.FieldsFunc(r.Message.To, split)

	// sender
	sender := r.Message.From
	if len(sender) == 0 {
		sender = r.Authentication.User
	}

	// Create mail
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", r.Message.Subject)

	// body
	if len(r.Message.Body) > 0 {
		m.SetBody("text/html", r.Message.Body)
	} else if len(r.Message.PlainBody) > 0 {
		m.SetBody("text/plain", r.Message.PlainBody)
		if len(r.Message.HtmlBody) > 0 {
			m.AddAlternative("text/html", r.Message.HtmlBody)
		}
	} else if len(r.Message.HtmlBody) > 0 {
		m.SetBody("text/html", r.Message.HtmlBody)
	}

	// smtp authentication
	d := gomail.NewDialer(r.Authentication.Server, r.Authentication.Port, r.Authentication.User, r.Authentication.Password)

	// Send mail
	return d.DialAndSend(m)
}

func split(r rune) bool {
	return r == ',' || r == ';'
}
