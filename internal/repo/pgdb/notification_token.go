package pgdb

import (
	"context"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	notificationTokensTable = "fcm_tokens"
)

type NotificationTokenRepo struct {
	*postgres.Database
}

func NewNotificationTokenRepo(db *postgres.Database) *NotificationTokenRepo {
	return &NotificationTokenRepo{db}
}

// SaveOrUpdate сохраняет новый токен или обновляет существующий
func (r *NotificationTokenRepo) SaveOrUpdate(
	ctx context.Context, log *slog.Logger, token entity.NotificationToken,
) error {
	log.Info("NotificationTokenRepo - SaveOrUpdate: started", "userID", token.UserID)

	sql, args, _ := r.Builder.Insert(notificationTokensTable).Columns(
		"user_id",
		"token",
	).Values(
		token.UserID,
		token.Token,
	).Suffix("ON CONFLICT (user_id) DO UPDATE SET token = EXCLUDED.token, updated_at = NOW()").ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("NotificationTokenRepo - SaveOrUpdate: failed", "error", err)
		return err
	}

	log.Info("NotificationTokenRepo - SaveOrUpdate: completed", "userID", token.UserID)
	return nil
}

// GetByUserID получает токен по userID
func (r *NotificationTokenRepo) GetByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) (*entity.NotificationToken, error) {
	log.Info("NotificationTokenRepo - GetByUserID: started", "userID", userID)

	sql, args, _ := r.Builder.Select(
		"id",
		"user_id",
		"token",
		"created_at",
		"updated_at",
	).From(notificationTokensTable).Where("user_id = ?", userID).ToSql()

	row := r.Cluster.QueryRow(ctx, sql, args...)

	var token entity.NotificationToken
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	if err != nil {
		log.Error("NotificationTokenRepo - GetByUserID: failed", "error", err)
		return nil, err
	}

	log.Info("NotificationTokenRepo - GetByUserID: completed", "userID", userID, "token", token.Token)
	return &token, nil
}
