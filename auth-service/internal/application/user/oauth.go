package user

import (
	"context"

	"github.com/LiliBeta/auth-service/internal/domain"
)

type OAuthService struct {
	userRepo    domain.UserRepository
	authService domain.AuthService
}

func NewOAuthService(userRepo domain.UserRepository, authService domain.AuthService) *OAuthService {
	return &OAuthService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *OAuthService) HandleOAuthCallback(ctx context.Context, provider, code string) (string, error) {
	// Implementar lógica para intercambiar código por token
	// Obtener información del usuario del proveedor OAuth
	// Crear/actualizar usuario en DB
	// Retornar token JWT

	return "", nil
}
