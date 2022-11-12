package controllers

import (
	"fmt"
	"github.com/gennesseaux/mailer/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	title := fmt.Sprintf("Mailer version %s", config.Version)
	secure := ""
	if config.C.Secure {
		secure = "(Running in secure mode)"
	}
	c.HTML(http.StatusOK, "html_index", gin.H{
		"title":  title,
		"secure": secure,
	})
}
