// repositories/notification_repository.go
package repositories

import (
	"context"
	"log"
	"notification-service/internal/entities"
)

type NotificationRepository interface {
	SendEmail(ctx context.Context, notification *entities.Notification) error
}

type SMTPNotificationRepository struct {
	smtpServer string
	fromEmail  string
	toEmail    string
}

func NewSMTPNotificationRepository(smtpServer, fromEmail, toEmail string) *SMTPNotificationRepository {
	return &SMTPNotificationRepository{
		smtpServer: smtpServer,
		fromEmail:  fromEmail,
		toEmail:    toEmail,
	}
}

func (r *SMTPNotificationRepository) SendEmail(ctx context.Context, notification *entities.Notification) error {
	// Implementación real de envío de email via SMTP
	// (simplificado para el ejemplo)
	log.Printf("Enviando email a %s: %+v", r.toEmail, notification)
	return nil
}
