package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type NotificationService struct {
	Client            *messaging.Client
	NotificationRepo  repo.Notification
	NotificationToken repo.NotificationToken
}

func NewNotificationService(
	ctx context.Context, app *firebase.App, notificationRepo repo.Notification,
	notificationTokenRepo repo.NotificationToken,
) *NotificationService {
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
	}
	return &NotificationService{
		Client: client, NotificationRepo: notificationRepo, NotificationToken: notificationTokenRepo,
	}
}

// SaveToken сохраняет токен в базу данных
func (n *NotificationService) SaveToken(ctx context.Context, log *slog.Logger, userID, token string) error {
	log.Info("NotificationService - SaveToken: started", "userID", userID, "token", token)

	err := n.NotificationToken.SaveOrUpdate(
		ctx, log, entity.NotificationToken{
			UserID: userID,
			Token:  token,
		},
	)
	if err != nil {
		log.Error("NotificationService - SaveToken: failed", "error", err)
		return err
	}

	log.Info("NotificationService - SaveToken: completed", "userID", userID)
	return nil
}

// getToken получает токен из базы данных
func (n *NotificationService) getToken(ctx context.Context, log *slog.Logger, userID string) (string, error) {
	log.Info("NotificationService - getToken: started", "userID", userID)

	token, err := n.NotificationToken.GetByUserID(ctx, log, userID)
	if err != nil {
		log.Error("NotificationService - getToken: failed", "error", err)
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	log.Info("NotificationService - getToken: completed", "userID", userID, "token", token.Token)
	return token.Token, nil
}

type NotificationCreateInput struct {
	UserID string
	Title  string
	Body   string
	Data   string
}

// SendNotification отправляет пуш-уведомление
func (n *NotificationService) SendNotification(
	ctx context.Context, log *slog.Logger, input NotificationCreateInput,
) error {
	// Получаем токен из базы данных
	token, err := n.getToken(ctx, log, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	log.Info("NotificationService - SendNotification: sending notification", "token", token)

	// Создаём сообщение для FCM
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: input.Title,
			Body:  input.Body,
		},
	}

	// Отправляем уведомление
	str, err := n.Client.Send(ctx, message)
	log.Info("NotificationService - SendNotification: result", "error", err, "response", str)
	if err != nil {
		return err
	}

	// Сохраняем уведомление в базе данных
	notification := entity.Notification{
		UserID: input.UserID,
		Title:  input.Title,
		Body:   input.Body,
		Data:   input.Data,
	}

	_, err = n.NotificationRepo.Create(ctx, notification)
	if err != nil {
		log.Error("NotificationService - SendNotification: failed to save notification", "error", err)
		return err
	}
	return nil
}

// GetNotificationsByUserID возвращает уведомления пользователя
func (n *NotificationService) GetNotificationsByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.Notification, error) {
	log.Info("NotificationService - GetNotificationsByUserID: started", "userID", userID)

	notifications, err := n.NotificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Error("NotificationService - GetNotificationsByUserID: failed", "error", err)
		return nil, err
	}

	log.Info("NotificationService - GetNotificationsByUserID: completed", "userID", userID)
	return notifications, nil
}
