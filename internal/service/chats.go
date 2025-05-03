package service

import (
	"context"
	"fmt"
	"log/slog"

	"family-flow-app/internal/entity"
	"family-flow-app/internal/repo"
)

type ChatMessageService struct {
	chatsRepo    repo.Chat
	messagesRepo repo.Message
}

func NewChatMessageService(chatsRepo repo.Chat, messagesRepo repo.Message) *ChatMessageService {
	return &ChatMessageService{
		chatsRepo:    chatsRepo,
		messagesRepo: messagesRepo,
	}
}

// Input для создания чата
type CreateChatInput struct {
	Name string
}

// Input для добавления участника в чат
type AddParticipantInput struct {
	ChatID string
	UserID string
}

// Input для создания чата с участниками
type CreateChatWithParticipantsInput struct {
	Name         string
	Participants []string
}

// Input для создания сообщения
type CreateMessageInput struct {
	ChatID   string
	SenderID string
	Content  string
}

// CreateChat создает новый чат
func (s *ChatMessageService) CreateChat(ctx context.Context, log *slog.Logger, input CreateChatInput) (string, error) {
	log.Info("Service - ChatMessageService - CreateChat")

	chat := entity.Chat{
		Name: input.Name,
	}

	id, err := s.chatsRepo.Create(ctx, chat)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - CreateChat: %v", err))
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	log.Info(fmt.Sprintf("Service - ChatMessageService - CreateChat - id: %s", id))
	return id, nil
}

// AddParticipant добавляет участника в чат
func (s *ChatMessageService) AddParticipant(ctx context.Context, log *slog.Logger, input AddParticipantInput) error {
	log.Info("Service - ChatMessageService - AddParticipant")

	err := s.chatsRepo.AddParticipant(ctx, input.ChatID, input.UserID)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - AddParticipant: %v", err))
		return fmt.Errorf("failed to add participant to chat: %w", err)
	}

	log.Info(
		fmt.Sprintf(
			"Service - ChatMessageService - AddParticipant - chatID: %s, userID: %s",
			input.ChatID,
			input.UserID,
		),
	)
	return nil
}

// CreateChatWithParticipants создает новый чат и добавляет участников
func (s *ChatMessageService) CreateChatWithParticipants(
	ctx context.Context, log *slog.Logger, input CreateChatWithParticipantsInput,
) (string, error) {
	log.Info("Service - ChatMessageService - CreateChatWithParticipants")

	// Создание чата
	chat := entity.Chat{
		Name: input.Name,
	}

	chatID, err := s.chatsRepo.Create(ctx, chat)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - CreateChatWithParticipants - CreateChat: %v", err))
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	// Добавление участников
	for _, userID := range input.Participants {
		err := s.chatsRepo.AddParticipant(ctx, chatID, userID)
		if err != nil {
			log.Error(
				fmt.Sprintf(
					"Service - ChatMessageService - CreateChatWithParticipants - AddParticipant: %v",
					err,
				),
			)
			return "", fmt.Errorf("failed to add participant to chat: %w", err)
		}
	}

	log.Info(fmt.Sprintf("Service - ChatMessageService - CreateChatWithParticipants - chatID: %s", chatID))
	return chatID, nil
}

// CreateMessage создает новое сообщение
func (s *ChatMessageService) CreateMessage(ctx context.Context, log *slog.Logger, input CreateMessageInput) (
	string, error,
) {
	log.Info("Service - ChatMessageService - CreateMessage")

	message := entity.Message{
		ChatID:   input.ChatID,
		SenderID: input.SenderID,
		Content:  input.Content,
	}

	id, err := s.messagesRepo.Create(ctx, message)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - CreateMessage: %v", err))
		return "", fmt.Errorf("failed to create message: %w", err)
	}

	log.Info(fmt.Sprintf("Service - ChatMessageService - CreateMessage - id: %s", id))
	return id, nil
}

// GetParticipants возвращает список участников чата
func (s *ChatMessageService) GetParticipants(ctx context.Context, log *slog.Logger, chatID string) ([]string, error) {
	log.Info("Service - ChatMessageService - GetParticipants")

	participants, err := s.chatsRepo.GetParticipants(ctx, chatID)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - GetParticipants: %v", err))
		return nil, fmt.Errorf("failed to get participants for chat: %w", err)
	}

	return participants, nil
}

// GetMessagesByChatID возвращает сообщения для указанного чата
func (s *ChatMessageService) GetMessagesByChatID(
	ctx context.Context, log *slog.Logger, chatID string,
) ([]entity.Message, error) {
	log.Info("Service - ChatMessageService - GetMessagesByChatID")

	messages, err := s.messagesRepo.GetByChatID(ctx, chatID)
	if err != nil {
		log.Error(fmt.Sprintf("Service - ChatMessageService - GetMessagesByChatID: %v", err))
		return nil, fmt.Errorf("failed to get messages by chat ID: %w", err)
	}

	return messages, nil
}
