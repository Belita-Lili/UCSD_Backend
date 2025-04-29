package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware establece un timeout para las peticiones
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
		}
	}
}
