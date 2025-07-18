// @title Subscriptions Service API
// @version 1.0
// @description REST-service for aggregating online user subscriptions
// @BasePath /
package main

import (
	"fmt"
	_ "github.com/Yagshymyradov/subscriptions-service/docs"
	"github.com/Yagshymyradov/subscriptions-service/internal/config"
	"github.com/Yagshymyradov/subscriptions-service/internal/db"
	"github.com/Yagshymyradov/subscriptions-service/internal/handlers"
	"github.com/Yagshymyradov/subscriptions-service/internal/repository"
	"github.com/Yagshymyradov/subscriptions-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
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
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.HTTP.Port)),
	))

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
