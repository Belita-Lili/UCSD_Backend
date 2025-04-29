package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "X-Request-ID"

// RequestIDMiddleware añade un ID único a cada petición
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Establecer en contexto y headers
		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set(requestIDKey, requestID)
		c.Next()
	}
}
