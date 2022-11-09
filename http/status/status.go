package httpstatus

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, status int, err error) {
	s := HttpStatus{
		Code:    status,
		Message: err.Error(),
	}
	c.JSON(status, s)
}

func Message(c *gin.Context, status int, obj interface{}) {
	s := HttpStatus{
		Code:    status,
		Message: obj,
	}
	c.JSON(status, s)
}

type HttpStatus struct {
	Code    int         `json:"code" example:"400"`
	Message interface{} `json:"message" example:"Bad Request"`
}
