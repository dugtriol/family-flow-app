package v1

import (
	"context"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"family-flow-app/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	chatsString = "/chats"
)

type ChatsRoutes struct {
	chatService service.Chats
}

func NewChatsRoutes(ctx context.Context, log *slog.Logger, route chi.Router, chatService service.Chats) {
	u := ChatsRoutes{chatService: chatService}
	route.Route(
		chatsString, func(r chi.Router) {
			r.Post("/", u.createChat(ctx, log))                                  // Создание чата
			r.Post("/{chatID}/participants", u.addParticipant(ctx, log))         // Добавление участника
			r.Post("/with-participants", u.createChatWithParticipants(ctx, log)) // Создание чата с участниками
			r.Get("/user", u.getChatsByUserID(ctx, log))                         // Получение чатов по ID пользователя
			r.Get("/{chatID}/messages", u.handleGetMessages(ctx, log))           // Получение сообщений по ID чата
		},
	)
}

type inputCreateChat struct {
	Name string `json:"name" validate:"required"`
}

type inputAddParticipant struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

// @Summary Create chat
// @Description Create a new chat
// @Tags chats
// @Accept json
// @Produce json
// @Param input body inputCreateChat true "Chat data"
// @Success 200 {string} string "Chat created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /chats [post]
func (u *ChatsRoutes) createChat(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handler - createChat - Start")

		var input inputCreateChat
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error("Handler - createChat - Failed to parse request", "error", err)
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			log.Error("Handler - createChat - Validation failed", "error", err)
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		chatID, err := u.chatService.CreateChat(ctx, log, service.CreateChatInput{Name: input.Name})
		if err != nil {
			log.Error("Handler - createChat - Failed to create chat", "error", err)
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create chat")
			return
		}

		log.Info("Handler - createChat - Chat created successfully", "chat_id", chatID)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"chat_id": chatID})
	}
}

// @Summary Add participant
// @Description Add a participant to a chat
// @Tags chats
// @Accept json
// @Produce json
// @Param chatID path string true "Chat ID"
// @Param input body inputAddParticipant true "Participant data"
// @Success 200 {string} string "Participant added"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /chats/{chatID}/participants [post]
func (u *ChatsRoutes) addParticipant(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handler - addParticipant - Start")

		chatID := chi.URLParam(r, "chatID")
		if chatID == "" {
			log.Error("Handler - addParticipant - Chat ID is required")
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Chat ID is required")
			return
		}

		var input inputAddParticipant
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error("Handler - addParticipant - Failed to parse request", "error", err)
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			log.Error("Handler - addParticipant - Validation failed", "error", err)
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		err := u.chatService.AddParticipant(
			ctx, log, service.AddParticipantInput{
				ChatID: chatID,
				UserID: input.UserID,
			},
		)
		if err != nil {
			log.Error("Handler - addParticipant - Failed to add participant", "error", err)
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to add participant")
			return
		}

		log.Info(
			"Handler - addParticipant - Participant added successfully",
			"chat_id",
			chatID,
			"user_id",
			input.UserID,
		)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Participant added successfully")
	}
}

type inputCreateChatWithParticipants struct {
	Name           string   `json:"name" validate:"required"`
	ParticipantIDs []string `json:"participant_ids" validate:"required"`
}

// create chat with participants
// @Summary Create chat with participants
// @Description Create a new chat with participants
// @Tags chats
// @Accept json
// @Produce json
// @Param input body inputCreateChatWithParticipants true "Chat data"
// @Success 200 {string} string "Chat created"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /chats/with-participants [post]
func (u *ChatsRoutes) createChatWithParticipants(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handler - createChatWithParticipants - Start")

		var input inputCreateChatWithParticipants

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			log.Error("Handler - createChatWithParticipants - Failed to get current user", "error", err)
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		// print r.Body
		log.Info("Handler - createChatWithParticipants - Request body", "body", r.Body)

		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error("Handler - createChatWithParticipants - Failed to parse request", "error", err)
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			log.Error("Handler - createChatWithParticipants - Validation failed", "error", err)
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		// log participant IDs
		log.Info("Handler - createChatWithParticipants - Participant IDs", "participant_ids", input.ParticipantIDs)

		input.ParticipantIDs = append(input.ParticipantIDs, user.Id)

		chatID, err := u.chatService.CreateChatWithParticipants(
			ctx, log, service.CreateChatWithParticipantsInput{
				Name:         input.Name,
				Participants: input.ParticipantIDs,
			},
		)
		if err != nil {
			log.Error("Handler - createChatWithParticipants - Failed to create chat", "error", err)
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create chat")
			return
		}

		log.Info("Handler - createChatWithParticipants - Chat created successfully", "chat_id", chatID)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]string{"chat_id": chatID})
	}
}

// get chats by user id
// @Summary Get chats by user ID
// @Description Get all chats for a user
// @Tags chats
// @Accept json
// @Produce json
// @Success 200 {array} entity.Chat "List of chats"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /chats/user [get]
// func (u *ChatsRoutes) getChatsByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Info("Handler - getChatsByUserID - Start")

// 		user, err := GetCurrentUserFromContext(r.Context())
// 		if err != nil {
// 			log.Error("Handler - getChatsByUserID - Failed to get current user", "error", err)
// 			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
// 			return
// 		}

// 		chats, err := u.chatService.GetChatsWithParticipants(ctx, log, user.Id)
// 		if err != nil {
// 			log.Error("Handler - getChatsByUserID - Failed to get chats", "error", err)
// 			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get chats")
// 			return
// 		}

//			log.Info("Handler - getChatsByUserID - Chats retrieved successfully", "user_id", user.Id)
//			w.WriteHeader(http.StatusOK)
//			render.JSON(w, r, chats)
//		}
//	}
//
// getChatsByUserID возвращает список чатов с последним сообщением
func (u *ChatsRoutes) getChatsByUserID(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handler - getChatsByUserID - Start")

		user, err := GetCurrentUserFromContext(r.Context())
		if err != nil {
			log.Error("Handler - getChatsByUserID - Failed to get current user", "error", err)
			response.NewError(w, r, log, err, http.StatusUnauthorized, "Failed to get current user")
			return
		}

		chats, err := u.chatService.GetChatsWithLastMessage(ctx, log, user.Id)
		if err != nil {
			log.Error("Handler - getChatsByUserID - Failed to get chats", "error", err)
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get chats")
			return
		}

		log.Info("Handler - getChatsByUserID - Chats retrieved successfully", "user_id", user.Id)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, chats)
	}
}

// @Summary Get messages by chat ID
// @Description Get all messages for a specific chat
// @Tags messages
// @Accept json
// @Produce json
// @Param chatID path string true "Chat ID"
// @Success 200 {array} entity.Message "List of messages"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /chats/{chatID}/messages [get]
func (u *ChatsRoutes) handleGetMessages(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handler - handleGetMessages - Start")

		// Получение chatID из параметров пути
		chatID := chi.URLParam(r, "chatID")
		if chatID == "" {
			log.Error("Handler - handleGetMessages - Chat ID is required")
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Chat ID is required")
			return
		}

		log.Info("Handler - handleGetMessages - Fetching messages", "chat_id", chatID)

		// Получение сообщений из сервиса
		messages, err := u.chatService.GetMessagesByChatID(ctx, log, chatID)
		if err != nil {
			log.Error("Handler - handleGetMessages - Failed to get messages", "error", err)
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to get messages")
			return
		}

		log.Info(
			"Handler - handleGetMessages - Messages retrieved successfully",
			"chat_id",
			chatID,
			"message_count",
			len(messages),
		)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, messages)
	}
}
