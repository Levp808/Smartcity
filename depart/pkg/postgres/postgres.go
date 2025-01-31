package postgres

import (
	"context"
	"depart/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func InitPostgresDbConnection(ctx context.Context, config *config.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	cfg, err := pgxpool.ParseConfig(databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	log.Info().Msg("Connected to postgres")

	return dbPool, nil
}
