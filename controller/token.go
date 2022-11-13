package controller

import (
	"github.com/gennesseaux/mailer/config"
	"github.com/gennesseaux/mailer/model"
	"github.com/gennesseaux/mailer/service/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tokenManager = auth.TokenManager{}

func Token(c *gin.Context) {
	var u model.User

	// Retrieve user login sent in the body
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid body provided")
		return
	}

	// Find user with username
	user := config.C.GetUser(u.Username)

	// Compare the user from the request, with the one defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	// Create an access token for the user
	td, err := tokenManager.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Send the token back
	tokens := map[string]string{
		"access_token": td.AccessToken,
	}

	c.JSON(http.StatusCreated, tokens)
}
