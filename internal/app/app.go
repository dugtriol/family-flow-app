package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"family-flow-app/config"
	v1 "family-flow-app/internal/handler/http/v1"
	"family-flow-app/internal/repo"
	"family-flow-app/internal/service"
	"family-flow-app/pkg/firebase"
	"family-flow-app/pkg/httpserver"
	"family-flow-app/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(cfg)

	//logger
	log := setLogger(cfg.Level)
	log.Info("Init logger")

	//rds, err := redis.New(ctx, log, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//log.Info("Init redis")

	//postgres
	database, err := postgres.New(ctx, cfg.Conn, postgres.MaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		fmt.Println(err.Error())
	}

	//repositories
	repos := repo.NewRepositories(database)
	ntf := firebase.Init(ctx)
	log.Info("Init firebase")
	dependencies := service.ServicesDependencies{
		Repos: repos, Config: cfg, App: ntf, BucketName: cfg.BucketName,
		Region: cfg.Region, EndpointResolver: cfg.EndpointResolver,
	}

	//services
	services := service.NewServices(ctx, dependencies)

	//handlers
	log.Info("Initializing handlers and routes...")

	router := chi.NewRouter()

	v1.NewRouter(ctx, log, router, services)
	// HTTP server
	log.Info("Starting http server...")
	log.Debug(fmt.Sprintf("Server port: %s", cfg.Port))
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
