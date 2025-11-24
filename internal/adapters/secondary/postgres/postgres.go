package postgres

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolConnection(config *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = config.Postgres.MaxConn
	poolConfig.MinConns = config.Postgres.MinConn

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
