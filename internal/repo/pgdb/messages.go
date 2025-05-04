package pgdb

import (
	"context"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"
)

const (
	messagesTable = "messages"
)

type MessagesRepo struct {
	*postgres.Database
}

func NewMessagesRepo(db *postgres.Database) *MessagesRepo {
	return &MessagesRepo{db}
}

// Create создает новое сообщение и возвращает его полностью
func (r *MessagesRepo) Create(ctx context.Context, message entity.Message) (entity.Message, error) {
	sql, args, _ := r.Builder.Insert(messagesTable).Columns(
		"chat_id",
		"sender_id",
		"content",
	).Values(
		message.ChatID,
		message.SenderID,
		message.Content,
	).Suffix("RETURNING id, chat_id, sender_id, content, created_at").ToSql()

	var createdMessage entity.Message
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(
		&createdMessage.ID,
		&createdMessage.ChatID,
		&createdMessage.SenderID,
		&createdMessage.Content,
		&createdMessage.CreatedAt,
	)
	if err != nil {
		return entity.Message{}, err
	}
	return createdMessage, nil
}

// GetByChatID возвращает сообщения для указанного чата
func (r *MessagesRepo) GetByChatID(ctx context.Context, chatID string) ([]entity.Message, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"chat_id",
		"sender_id",
		"content",
		"created_at",
	).From(messagesTable).Where(
		"chat_id = ?", chatID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var message entity.Message
		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
