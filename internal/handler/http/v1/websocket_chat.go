package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все соединения (для тестов)
	},
}

type WebSocketRequest struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type WebSocketResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WebSocketHandler(ctx context.Context, log *slog.Logger, chatService service.Chats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Info(fmt.Sprintf("Failed to upgrade connection:", err))
			return
		}
		defer conn.Close()

		for {
			// Чтение сообщения от клиента
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Info("Error reading message:", err)
				break
			}

			// Обработка сообщения
			var req WebSocketRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Info("Invalid message format:", err)
				sendError(conn, "Invalid message format")
				continue
			}

			// Выполнение действия
			switch req.Action {
			case "send_message":
				handleSendMessage(ctx, log, conn, chatService, req.Data)
			case "get_messages":
				handleGetMessages(ctx, log, conn, chatService, req.Data)
			default:
				sendError(conn, "Unknown action")
			}
		}
	}
}

func handleGetMessages(
	ctx context.Context, log *slog.Logger, conn *websocket.Conn, chatService service.Chats, data json.RawMessage,
) {
	var input struct {
		ChatID string `json:"chat_id"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		sendError(conn, "Invalid input for get_messages")
		return
	}

	messages, err := chatService.GetMessagesByChatID(ctx, log, input.ChatID)
	if err != nil {
		sendError(conn, "Failed to get messages: "+err.Error())
		return
	}

	sendSuccess(conn, messages)
}

func sendError(conn *websocket.Conn, message string) {
	resp := WebSocketResponse{
		Status:  "error",
		Message: message,
	}
	conn.WriteJSON(resp)
}

func sendSuccess(conn *websocket.Conn, data interface{}) {
	resp := WebSocketResponse{
		Status: "success",
		Data:   data,
	}
	conn.WriteJSON(resp)
}

func handleSendMessage(
	ctx context.Context, log *slog.Logger, conn *websocket.Conn, chatService service.Chats, data json.RawMessage,
) {
	var input service.CreateMessageInput
	if err := json.Unmarshal(data, &input); err != nil {
		sendError(conn, "Invalid input for send_message")
		return
	}

	messageID, err := chatService.CreateMessage(ctx, log, input)
	if err != nil {
		sendError(conn, "Failed to send message: "+err.Error())
		return
	}

	sendSuccess(conn, map[string]string{"message_id": messageID})
}
