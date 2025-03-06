package db

import (
	"context"
	"log"
	"time"
	"tritchgo/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	defaultMaxConns          = 5
	defaultMinConns          = 1
	defaultMaxConnLifetime   = 30 * time.Minute
	defaultMaxConnIdleTime   = 15 * time.Minute
	defaultHealthCheckPeriod = 2 * time.Minute
	defaultConnectTimeout    = 5 * time.Second
)

func DBConn(ctx context.Context) (*pgxpool.Pool, error) {
	url := config.LoadEnv().DB_URL

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to parse config DBPOOL: %v", err)
	}

	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.MaxConnLifetime = defaultMaxConnLifetime
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.HealthCheckPeriod = defaultHealthCheckPeriod
	config.ConnConfig.ConnectTimeout = defaultConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database pool created successfully")
	return pool, nil
}
