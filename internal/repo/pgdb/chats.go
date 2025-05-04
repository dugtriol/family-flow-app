package pgdb

import (
	"context"

	"family-flow-app/internal/entity"
	"family-flow-app/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	chatsTable            = "chats"
	chatParticipantsTable = "chat_participants"
)

type ChatsRepo struct {
	*postgres.Database
}

func NewChatsRepo(db *postgres.Database) *ChatsRepo {
	return &ChatsRepo{db}
}

// Create создает новый чат
func (r *ChatsRepo) Create(ctx context.Context, chat entity.Chat) (string, error) {
	sql, args, _ := r.Builder.Insert(chatsTable).Columns(
		"name",
	).Values(
		chat.Name,
	).Suffix("RETURNING id").ToSql()

	var id string
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// AddParticipant добавляет участника в чат
func (r *ChatsRepo) AddParticipant(ctx context.Context, chatID, userID string) error {
	sql, args, _ := r.Builder.Insert(chatParticipantsTable).Columns(
		"chat_id",
		"user_id",
	).Values(
		chatID,
		userID,
	).ToSql()

	_, err := r.Cluster.Exec(ctx, sql, args...)
	return err
}

// GetParticipants возвращает список участников чата
func (r *ChatsRepo) GetParticipants(ctx context.Context, chatID string) ([]string, error) {
	sql, args, _ := r.Builder.Select(
		"user_id",
	).From(chatParticipantsTable).Where(
		"chat_id = ?", chatID,
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []string
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}
	return participants, nil
}

// GetByID возвращает чат по его ID
func (r *ChatsRepo) GetByID(ctx context.Context, id string) (entity.Chat, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"name",
		"created_at",
	).From(chatsTable).Where(
		"id = ?", id,
	).ToSql()

	var chat entity.Chat
	err := r.Cluster.QueryRow(ctx, sql, args...).Scan(
		&chat.ID,
		&chat.Name,
		&chat.CreatedAt,
	)
	if err != nil {
		return entity.Chat{}, err
	}
	return chat, nil
}

// GetAll возвращает список всех чатов
func (r *ChatsRepo) GetAll(ctx context.Context) ([]entity.Chat, error) {
	sql, args, _ := r.Builder.Select(
		"id",
		"name",
		"created_at",
	).From(chatsTable).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

// get chats with participants
func (r *ChatsRepo) GetChatsWithParticipants(ctx context.Context, userID string) ([]entity.Chat, error) {
	sql, args, _ := r.Builder.Select(
		"c.id",
		"c.name",
		"c.created_at",
	).From(chatsTable + " c").Join(
		chatParticipantsTable + " cp ON c.id = cp.chat_id", // Добавлено условие соединения
	).Where(
		squirrel.Eq{"cp.user_id": userID},
	).ToSql()

	rows, err := r.Cluster.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
