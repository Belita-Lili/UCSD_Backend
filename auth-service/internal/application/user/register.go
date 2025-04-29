// internal/application/user/register.go
package user

import (
	"context"
	"errors"

	"github.com/LiliBeta/auth-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	userRepo domain.UserRepository
}

func NewRegisterService(userRepo domain.UserRepository) *RegisterService {
	return &RegisterService{
		userRepo: userRepo,
	}
}

func (s *RegisterService) RegisterUser(ctx context.Context, user domain.User) error {
	// Validar que el email no exista
	existing, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrEmailAlreadyExists
	}

	// Hash de contraseña con bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Guardar usuario
	return s.userRepo.Save(ctx, &user)
}
