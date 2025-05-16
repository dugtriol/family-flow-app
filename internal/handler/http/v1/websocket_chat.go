package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"family-flow-app/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все соединения (для тестов)
	},
}

// Хранилище для активных WebSocket-соединений
var connections = struct {
	sync.Mutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
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

func sendError(conn *websocket.Conn, message string) {
	log := slog.Default()
	log.Error("sendError - Sending error response", "message", message)

	resp := WebSocketResponse{
		Status:  "error",
		Message: message,
	}
	conn.WriteJSON(resp)
}

func WebSocketHandler(ctx context.Context, log *slog.Logger, chatService service.Chats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("WebSocketHandler - Start connection upgrade")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("WebSocketHandler - Failed to upgrade connection", "error", err)
			return
		}
		defer func() {
			connections.Lock()
			delete(connections.clients, conn)
			connections.Unlock()
			conn.Close()
			log.Info("WebSocketHandler - Connection closed")
		}()

		// Добавляем соединение в хранилище
		connections.Lock()
		connections.clients[conn] = true
		connections.Unlock()
		log.Info("WebSocketHandler - Connection upgraded successfully")

		// Устанавливаем pong handler для продления соединения
		conn.SetPongHandler(
			func(appData string) error {
				log.Info("WebSocketHandler - Pong received")
				return nil
			},
		)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("WebSocketHandler - Error reading message", "error", err)
				break
			}

			log.Info("WebSocketHandler - Message received", "message", string(message))

			// Обработка ping от клиента
			var raw map[string]interface{}
			if err := json.Unmarshal(message, &raw); err == nil {
				if t, ok := raw["type"]; ok && t == "ping" {
					// Отправляем pong в ответ
					conn.WriteJSON(map[string]string{"type": "pong"})
					log.Info("WebSocketHandler - Pong sent")
					continue
				}
			}

			// Обработка остальных сообщений
			var req WebSocketRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Error("WebSocketHandler - Invalid message format", "error", err)
				sendError(conn, "Invalid message format")
				continue
			}

			log.Info("WebSocketHandler - Processing action", "action", req.Action)

			switch req.Action {
			case "send_message":
				handleSendMessage(ctx, log, conn, chatService, req.Data)
			default:
				log.Warn("WebSocketHandler - Unknown action", "action", req.Action)
				sendError(conn, "Unknown action")
			}
		}
	}
}

type CreateMessageInput struct {
	ChatID   string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
}

func handleSendMessage(
	ctx context.Context, log *slog.Logger, conn *websocket.Conn, chatService service.Chats, data json.RawMessage,
) {
	log.Info("handleSendMessage - Start")

	var input CreateMessageInput
	if err := json.Unmarshal(data, &input); err != nil {
		log.Error("handleSendMessage - Invalid input", "error", err)
		sendError(conn, "Invalid input for send_message")
		return
	}

	log.Info("handleSendMessage - Sending message", "chat_id", input.ChatID, "sender_id", input.SenderID)

	output, err := chatService.CreateMessage(
		ctx, log, service.CreateMessageInput{
			ChatID:   input.ChatID,
			SenderID: input.SenderID,
			Content:  input.Content,
		},
	)
	if err != nil {
		log.Error("handleSendMessage - Failed to send message", "error", err)
		sendError(conn, "Failed to send message: "+err.Error())
		return
	}

	log.Info("handleSendMessage - Message sent successfully", "message:", output)
	connections.Lock()
	log.Info("handleSendMessage - Broadcasting message", "connected_clients", len(connections.clients))
	defer connections.Unlock()
	for client := range connections.clients {
		if client == conn {
			continue // Не отправляем сообщение обратно тому, кто его отправил
		}
		if err := client.WriteJSON(
			WebSocketResponse{
				Status: "success",
				Data:   output,
			},
		); err != nil {
			log.Error("handleSendMessage - Failed to send message to client", "error", err)
			client.Close()
			delete(connections.clients, client)
		}
	}
}
