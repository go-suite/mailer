package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-suite/mailer/config"
	"net/http"
)

type info struct {
	Version string
	Secure  bool
}

func Info(c *gin.Context) {
	i := info{
		Version: config.Version,
		Secure:  config.C.Secure,
	}
	c.JSON(http.StatusOK, i)
}
