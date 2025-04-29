package middlewares

import (
	"net/http"

	"github.com/LiliBeta/auth-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService domain.AuthService
}

func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		return
	}

	valid, err := m.authService.VerifyToken(c.Request.Context(), token)
	if err != nil || !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.Next()
}
