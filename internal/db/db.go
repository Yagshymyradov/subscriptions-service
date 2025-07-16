package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/Yagshymyradov/subscriptions-service/internal/config"
	"github.com/jackc/pgx/v5/stdlib"
)

func New(cfg *config.DBConfig) (*sql.DB, error) {
	sql.Register("pqx", stdlib.GetDefaultDriver())

	db, err := sql.Open("pqx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}