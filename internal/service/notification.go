package service

import (
	"context"
	`fmt`
	"log"
	`log/slog`

	`family-flow-app/internal/entity`
	`family-flow-app/internal/repo`
	`family-flow-app/pkg/redis`
	firebase `firebase.google.com/go/v4`
	"firebase.google.com/go/v4/messaging"
)

type NotificationService struct {
	Rd               *redis.Redis
	Client           *messaging.Client
	NotificationRepo repo.Notification
}

func NewNotificationService(
	ctx context.Context, rd *redis.Redis, app *firebase.App, notificationRepo repo.Notification,
) *NotificationService {
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
	}
	return &NotificationService{Rd: rd, Client: client, NotificationRepo: notificationRepo}
}

func (n *NotificationService) SaveToken(ctx context.Context, userID, token string) error {
	// Save the token in Redis with the user ID as the key
	statusCmd := n.Rd.Set(ctx, userID, token, 0) // 0 means no expiration
	if statusCmd.Err() != nil {
		return fmt.Errorf("failed to save token in Redis: %w", statusCmd.Err())
	}
	return nil
}

func (n *NotificationService) getToken(ctx context.Context, userID string) (string, error) {
	statusCmd := n.Rd.Get(ctx, userID)
	if statusCmd.Err() != nil {
		return "", fmt.Errorf("failed to get token from Redis: %w", statusCmd.Err())
	}
	return statusCmd.Val(), nil
}

type NotificationCreateInput struct {
	UserID string
	Title  string
	Body   string
	Data   string
}

func (n *NotificationService) SendNotification(
	ctx context.Context, log *slog.Logger, input NotificationCreateInput,
) error {
	var err error
	// Retrieve the token from Redis
	statusCmd := n.Rd.Get(ctx, input.UserID)
	if statusCmd.Err() != nil {
		return fmt.Errorf("failed to get token from Redis: %w", statusCmd.Err())
	}

	token := statusCmd.Val()

	// Create the FCM message
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: input.Title,
			Body:  input.Body,
		},
	}

	// Send the notification
	str, err := n.Client.Send(ctx, message)
	log.Info(fmt.Sprintf("Service - NotificationService - SendNotification: %v %s", err, str))
	if err != nil {
		return err
	}

	// Save the notification in the database
	notification := entity.Notification{
		UserID: input.UserID,
		Title:  input.Title,
		Body:   input.Body,
		Data:   input.Data,
	}

	_, err = n.NotificationRepo.Create(ctx, notification)
	if err != nil {
		log.Error("Service - NotificationService - SendNotification: %v", err)
		return err
	}
	return nil
}

func (n *NotificationService) GetNotificationsByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.Notification, error) {
	notifications, err := n.NotificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Error("Service - NotificationService - GetNotificationsByUserID: %v", err)
		return nil, err
	}
	return notifications, nil
}
