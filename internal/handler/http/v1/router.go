package v1

import (
	"context"
	"log/slog"
	"net/http"

	"family-flow-app/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	mwLogger "family-flow-app/pkg/middleware"
)

const api = "/api"

func NewRouter(ctx context.Context, log *slog.Logger, route *chi.Mux, services *service.Services) {
	route.Use(middleware.Logger)
	route.Use(middleware.RequestID)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mwLogger.New(log))
	route.Use(render.SetContentType(render.ContentTypeJSON))

	route.Route(
		api, func(r chi.Router) {
			r.Get("/ping", Ping())
			NewAuthRoutes(ctx, log, r, services.User)
			NewEmailRoutes(ctx, log, r, services.Email)
			r.Group(
				func(g chi.Router) {
					g.Use(AuthMiddleware(ctx, log, services.User))
					NewUserRoutes(ctx, log, g, services.User)
					NewTaskRoutes(ctx, log, g, services.Task)
					NewFamilyRoutes(ctx, log, g, services.Email, services.Family)
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
