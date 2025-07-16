package main

import (
	"github.com/Yagshymyradov/subscriptions-service/internal/config"
	"github.com/Yagshymyradov/subscriptions-service/internal/db"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main(){
	_ = godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	logger.Info("config DSN", zap.String("dsn", cfg.DB.DSN))

	logger.Info("Starting service", zap.String("port", cfg.HTTP.Port))

	_, err = db.New(&cfg.DB)
	if err != nil {
		logger.Fatal("db connection failed", zap.Error(err))
	}

	logger.Info("Connected to database")
}