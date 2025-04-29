package user

import (
	"context"

	"github.com/LiliBeta/auth-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo    domain.UserRepository
	authService domain.AuthService
}

func NewLoginService(userRepo domain.UserRepository, authService domain.AuthService) *LoginService {
	return &LoginService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *LoginService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.authService.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	return token, nil
}
