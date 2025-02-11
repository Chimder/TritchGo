package sqlc

import (
	"context"
	"log"
	"tritchgo/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func DBConn(ctx context.Context) (*pgxpool.Pool, error) {
	url := config.LoadEnv().DB_URL
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
		return nil, err
	}

	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	sqlcPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil, err
	}

	return sqlcPool, nil
}
