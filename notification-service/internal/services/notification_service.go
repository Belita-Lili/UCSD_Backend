// services/notification_service.go
package services

import (
	"context"
	"notification-service/internal/entities"
	"notification-service/internal/repositories"
	"time"
)

type NotificationService interface {
	SendErrorNotification(ctx context.Context, errorMsg string, metadata string) error
}

type notificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) *notificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) SendErrorNotification(ctx context.Context, errorMsg string, metadata string) error {
	notification := &entities.Notification{
		Error:     errorMsg,
		Timestamp: time.Now().Format(time.RFC3339),
		Metadata:  metadata,
	}

	return s.repo.SendEmail(ctx, notification)
}
