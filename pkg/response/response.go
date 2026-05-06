package response

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func OK(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, message string, errors any) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
