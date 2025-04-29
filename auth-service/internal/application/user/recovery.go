package user

import (
	"context"
	"time"

	"github.com/LiliBeta/auth-service/internal/domain"
	"github.com/google/uuid"
)

type RecoveryService struct {
	userRepo domain.UserRepository
	emailCfg EmailConfig
}

type EmailConfig struct {
	From     string
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

func NewRecoveryService(userRepo domain.UserRepository, emailCfg EmailConfig) *RecoveryService {
	return &RecoveryService{
		userRepo: userRepo,
		emailCfg: emailCfg,
	}
}

func (s *RecoveryService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 1)

	// Guardar token en DB (implementar método en repositorio)
	if err := s.userRepo.SaveResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return err
	}

	// Enviar email (implementar servicio de email)
	resetLink := "https://yourapp.com/reset-password?token=" + token
	emailBody := "Click the link to reset your password: " + resetLink

	// email.Send(s.emailCfg, email, "Password Reset", emailBody)

	return nil
}

func (s *RecoveryService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validar token y cambiar contraseña
	// Implementar lógica completa
	return nil
}
