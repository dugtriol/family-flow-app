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
			r.Post("/", u.createChat(ctx, log))                          // Создание чата
			r.Post("/{chatID}/participants", u.addParticipant(ctx, log)) // Добавление участника
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
		var input inputCreateChat
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		chatID, err := u.chatService.CreateChat(ctx, log, service.CreateChatInput{Name: input.Name})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to create chat")
			return
		}

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
		chatID := chi.URLParam(r, "chatID")
		if chatID == "" {
			response.NewError(w, r, log, nil, http.StatusBadRequest, "Chat ID is required")
			return
		}

		var input inputAddParticipant
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			response.NewError(w, r, log, err, http.StatusBadRequest, "Failed to parse request")
			return
		}
		if err := validator.New().Struct(input); err != nil {
			response.NewValidateError(w, r, log, http.StatusBadRequest, "Invalid request", err)
			return
		}

		err := u.chatService.AddParticipant(ctx, log, service.AddParticipantInput{
			ChatID: chatID,
			UserID: input.UserID,
		})
		if err != nil {
			response.NewError(w, r, log, err, http.StatusInternalServerError, "Failed to add participant")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "Participant added successfully")
	}
}
