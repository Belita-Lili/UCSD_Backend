package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: status >= http.StatusOK && status < http.StatusMultipleChoices,
		Data:    data,
	})
}

func RespondError(c *gin.Context, status int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}
