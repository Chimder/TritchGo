package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepo(db *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{db: db}
}

func (st *StatsRepo) GetUserStatsById(ctx context.Context, id string) ([]StreamStats, error) {
	query := `SELECT * FROM stream_stats WHERE user_id = $1`
	rows, err := st.db.Query(ctx, query, id)
	if err != nil {
		log.Printf("Err fetch user stats  %v", err)
		return nil, err
	}
	defer rows.Close()

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[StreamStats])

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
	}
	return stats, err
}

func (st *StatsRepo) GetStreamStatsById(ctx context.Context, id string) ([]StreamStats, error) {
	query := `SELECT * FROM stream_stats WHERE stream_id = $1`

	rows, err := st.db.Query(ctx, query, id)
	if err != nil {
		log.Printf("Err fetch stream stats  %v", err)
		return nil, err
	}

	s, err := pgx.CollectRows(rows, pgx.RowToStructByName[StreamStats])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning row: %v", err)
		return nil, err
	}
	return s, nil
}
