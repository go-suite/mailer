package controller

import (
	"github.com/gin-gonic/gin"
	status "github.com/go-suite/mailer/http/status"
	"net/http"
)

func Check(c *gin.Context) {
	status.Message(c, http.StatusOK, "Ok")
}
