package controllers

import (
	"github.com/gennesseaux/mailer/config"
	status "github.com/gennesseaux/mailer/http/status"
	"github.com/gennesseaux/mailer/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"net/http"
	"strings"
)

func Send(c *gin.Context) {
	var r models.Request
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

func sendMail(r models.Request) error {
	tos := strings.FieldsFunc(r.Message.To, split)

	m := gomail.NewMessage()
	m.SetHeader("From", r.Message.From)
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", r.Message.Subject)
	m.SetBody("text/html", r.Message.Body)

	d := gomail.NewDialer(r.Authentication.Server, r.Authentication.Port, r.Authentication.User, r.Authentication.Password)

	return d.DialAndSend(m)
}

func split(r rune) bool {
	return r == ',' || r == ';'
}
