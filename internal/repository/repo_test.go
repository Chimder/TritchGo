package repository

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func initDb(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:latest",
		tcpostgres.WithDatabase("test"),
		tcpostgres.WithUsername("user"),
		tcpostgres.WithPassword("password"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	sqlDB, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	err = goose.Up(sqlDB, "../../sql/migration")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	return pool
}

func GetSteam() *StreamArg {

	stream := StreamArg{
		ID:          gofakeit.UUID(),
		UserID:      gofakeit.DigitN(9),
		Title:       gofakeit.Sentence(),
		GameID:      gofakeit.DigitN(6),
		ViewerCount: gofakeit.IntRange(0, 100000),
	}
	return &stream
}
func TestRepo(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	gofakeit.Seed(time.Now().UnixNano())
	dbPool := initDb(t)
	require.NotNil(t, dbPool)
	statsRepo := NewStatsRepo(dbPool)

	var batchMutex sync.Mutex
	var wg sync.WaitGroup
	streamBatchInsert := &pgx.Batch{}

	for i := range 5000 {
		wg.Add(1)
		startedAt := time.Now().Add(-time.Duration(gofakeit.Number(1, 720)) * time.Minute)

		go func(i int) {
			defer wg.Done()

			queryStr, pgxArg := statsRepo.BuildInsertQuery(GetSteam(), startedAt)
			batchMutex.Lock()
			streamBatchInsert.Queue(queryStr, pgxArg)
			batchMutex.Unlock()
		}(i)

	}
	wg.Wait()

	res := dbPool.SendBatch(ctx, streamBatchInsert)
	defer res.Close()

	for i := 0; i < streamBatchInsert.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			assert.NoError(t, err, "failed to execute batch insert")
		}
	}

	t.Run("Get all stats", func(t *testing.T) {
		allStats, err := statsRepo.GetAllStats()
		assert.NotEmpty(t, allStats)
		assert.NoError(t, err)

		for _, v := range allStats {

			/////////////////////// get stream by id
			streamById, err := statsRepo.GetStreamStatsById(t.Context(), v.StreamID)
			assert.NotEmpty(t, streamById)
			assert.NoError(t, err)

			////////////////////////////// get steam by user id
			userId, err := statsRepo.GetUserStatsById(t.Context(), v.UserID)
			assert.NotEmpty(t, userId)
			assert.NoError(t, err)
		}
	})
	// time.Sleep(5 * time.Minute)
}
