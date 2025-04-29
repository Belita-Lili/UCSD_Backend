package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			switch e := err.Err.(type) {
			case *AppError:
				c.JSON(e.Code, gin.H{
					"error":   e.Message,
					"details": e.Err,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}
	}
}
