package v1

import (
	"context"
	"log/slog"
	"net/http"

	_ "family-flow-app/docs"
	"family-flow-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	mwLogger "family-flow-app/pkg/middleware"
)

const api = "/api"

func NewRouter(ctx context.Context, log *slog.Logger, route *chi.Mux, services *service.Services) {
	log.Info("Initializing router")
	route.Use(middleware.Logger)
	route.Use(middleware.RequestID)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mwLogger.New(log))
	route.Use(render.SetContentType(render.ContentTypeJSON))

	log.Info("Initializing websocket..")
	route.Group(
		func(r chi.Router) {
			r.HandleFunc(
				"/ws", WebSocketHandler(ctx, log, services.Chats),
			)
		},
	)

	log.Info("Swagger is available")
	route.Get(
		"/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		),
	)

	route.Route(
		api, func(r chi.Router) {
			r.Get("/ping", Ping())
			NewAuthRoutes(ctx, log, r, services.User)
			NewEmailRoutes(ctx, log, r, services.Email)
			r.Group(
				func(g chi.Router) {
					g.Use(AuthMiddleware(ctx, log, services.User))
					NewUserRoutes(ctx, log, g, services.User, services.File)
					NewFamilyRoutes(ctx, log, g, services.Email, services.Family, services.File, services.Notification)
					NewTodoRoutes(ctx, log, g, services.TodoItem, services.Notification)
					NewShoppingRoutes(ctx, log, g, services.ShoppingItem, services.Notification, services.Family)
					NewWishlistRoutes(ctx, log, g, services.WishlistItem, services.Notification, services.Family)
					NewNotificationRoutes(ctx, log, g, services.Notification)
					NewChatsRoutes(ctx, log, g, services.Chats)
					NewRewardsRoutes(ctx, log, g, services.Rewards, services.Notification, services.Family)
					NewDiaryRoutes(ctx, log, g, services.Diary)
				},
			)
		},
	)
}

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ok"))

	}
}
