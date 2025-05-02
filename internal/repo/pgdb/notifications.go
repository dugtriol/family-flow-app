package pgdb

import (
	`context`

	`family-flow-app/internal/entity`
	`family-flow-app/pkg/postgres`
)

const (
	notificationsTable = "notifications"
)

type NotificationsRepo struct {
	*postgres.Database
}

func NewNotificationsRepo(db *postgres.Database) *NotificationsRepo {
	return &NotificationsRepo{db}
}

func (r *NotificationsRepo) Create(ctx context.Context, notification entity.Notification) (string, error) {
	sql, args, _ := r.Builder.Insert(notificationsTable).Columns(
		"user_id",
		"title",
		"body",
		"data",
	).Values(
		notification.UserID,
		notification.Title,
		notification.Body,
		notification.Data,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *NotificationsRepo) GetByUserID(ctx context.Context, userID string) ([]entity.Notification, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"user_id",
		"title",
		"body",
		"data",
		"is_read",
		"created_at",
	).From(notificationsTable).Where(
		"user_id = ?", userID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []entity.Notification
	for rows.Next() {
		var notification entity.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Body,
			&notification.Data,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
