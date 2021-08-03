package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func ConnectPool(dbURL string, logger zerolog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	config.ConnConfig.Logger = zerologadapter.NewLogger(logger)

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to poll: %w", err)
	}

	return db, nil
}
