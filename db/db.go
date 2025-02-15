package db

import (
	"context"
	"log"
	"tritchgo/config"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func DBConn(ctx context.Context) (*pgx.Conn, error) {
	url := config.LoadEnv().DB_URL
	// cfg, err := pgxpool.ParseConfig(url)
	// if err != nil {
	// 	log.Fatalf("Unable to parse config: %v", err)
	// 	return nil, err
	// }

	// cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	db, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil, err
	}

	return db, nil
}
