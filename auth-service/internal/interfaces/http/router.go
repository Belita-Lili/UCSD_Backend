package http

import (
	"github.com/LiliBeta/auth-service/internal/interfaces/http/handlers"
	"github.com/LiliBeta/auth-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userHandlers *handlers.UserHandler,
	authHandlers *handlers.AuthHandler,
	authMiddleware *middlewares.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// Middlewares globales
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.SecurityHeadersMiddleware())

	// Grupo de rutas públicas
	public := router.Group("/api/v1/auth")
	{
		public.POST("/register", userHandlers.Register)
		public.POST("/login", userHandlers.Login)
		public.POST("/password-reset", userHandlers.RequestPasswordReset)
		public.POST("/password-reset/confirm", userHandlers.ResetPassword)
		public.GET("/oauth/:provider", userHandlers.OAuthRedirect)
		public.GET("/oauth/:provider/callback", userHandlers.OAuthCallback)
	}

	// Grupo de rutas protegidas
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.ValidateToken)
	{
		// Rutas protegidas...
	}

	return router
}
