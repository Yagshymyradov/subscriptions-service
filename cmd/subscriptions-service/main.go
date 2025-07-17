package main

import (
	"net/http"

	"github.com/Yagshymyradov/subscriptions-service/internal/config"
	"github.com/Yagshymyradov/subscriptions-service/internal/db"
	"github.com/Yagshymyradov/subscriptions-service/internal/handlers"
	"github.com/Yagshymyradov/subscriptions-service/internal/repository"
	"github.com/Yagshymyradov/subscriptions-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	database, err := db.New(&cfg.DB)
	if err != nil {
		logger.Fatal("db connection failed", zap.Error(err))
	}
	logger.Info("Connected to database")

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	subRepo := repository.NewPostgresSubscriptionRepository(database)
	subSvc := service.NewSubscriptionService(subRepo)
	h := handlers.New(subSvc, logger)
	handlers.RegisterRoutes(router, h)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: router,
	}
	logger.Info("http listen", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
}
