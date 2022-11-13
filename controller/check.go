package controller

import (
	status "github.com/gennesseaux/mailer/http/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Check(c *gin.Context) {
	status.Message(c, http.StatusOK, "Ok")
}
