package repository

// import (
// 	"context"
// 	"database/sql"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	"github.com/pressly/goose/v3"
// 	"github.com/stretchr/testify/require"
// 	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
// )

// func initDb(t *testing.T) *pgxpool.Pool {
// 	ctx := context.Background()
// 	pgContainer, err := tcpostgres.Run(ctx,
// 		"postgres:latest",
// 		tcpostgres.WithDatabase("test"),
// 		tcpostgres.WithUsername("user"),
// 		tcpostgres.WithPassword("password"),
// 		tcpostgres.BasicWaitStrategies(),
// 	)
// 	require.NoError(t, err)

// 	t.Cleanup(func() {
// 		_ = pgContainer.Terminate(ctx)
// 	})

// 	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
// 	require.NoError(t, err)

// 	sqlDB, err := sql.Open("pgx", connStr)
// 	require.NoError(t, err)

// 	err = goose.Up(sqlDB, "../../sql/migration")
// 	require.NoError(t, err)

// 	pool, err := pgxpool.New(ctx, connStr)
// 	require.NoError(t, err)

// 	return pool
// }
// func TestRepo(t *testing.T) {
// 	// dbPool := initDb(t)
// 	// statsRepo := NewStatsRepo(dbPool)

// 	// statsRepo
// }
