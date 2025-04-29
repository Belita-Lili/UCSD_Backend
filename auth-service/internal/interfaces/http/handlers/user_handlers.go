package handlers

import (
	"net/http"

	"github.com/LiliBeta/auth-service/internal/application/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	registerService *user.RegisterService
	loginService    *user.LoginService
	recoveryService *user.RecoveryService
	oauthService    *user.OAuthService
}

func NewUserHandler(
	registerService *user.RegisterService,
	loginService *user.LoginService,
	recoveryService *user.RecoveryService,
	oauthService *user.OAuthService,
) *UserHandler {
	return &UserHandler{
		registerService: registerService,
		loginService:    loginService,
		recoveryService: recoveryService,
		oauthService:    oauthService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	// Implementación ya mostrada
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.loginService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.recoveryService.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}

func (h *UserHandler) OAuthRedirect(c *gin.Context) {
	provider := c.Param("provider")
	state := uuid.New().String()

	// Guardar state en sesión o cache para validación posterior

	var authURL string
	switch provider {
	case "google":
		authURL = h.oauthService.GetGoogleAuthURL(state)
	case "facebook":
		authURL = h.oauthService.GetFacebookAuthURL(state)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *UserHandler) OAuthCallback(c *gin.Context) {
	// Implementar manejo de callback OAuth
}
