package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepo interface {
	GetUserStatsById(ctx context.Context, id string) ([]StreamStats, error)
	GetStreamStatsById(ctx context.Context, id string) ([]StreamStats, error)
}

type statsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepo(db *pgxpool.Pool) StatsRepo {
	return &statsRepo{db: db}
}

func (r *statsRepo) GetUserStatsById(ctx context.Context, id string) ([]StreamStats, error) {
	query := `SELECT * FROM stream_stats WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("err fetch user stats  %w", err)
	}
	defer rows.Close()

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[StreamStats])

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("err iterating rows %w", err)
	}

	return stats, err
}

func (r *statsRepo) GetStreamStatsById(ctx context.Context, id string) ([]StreamStats, error) {

	rows, err := r.db.Query(ctx, `SELECT * FROM stream_stats WHERE stream_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("err fetch stream stats %w", err)
	}

	s, err := pgx.CollectRows(rows, pgx.RowToStructByName[StreamStats])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no stream rows found: %w", err)
		}
		return nil, fmt.Errorf("err scanning row %w", err)
	}
	return s, nil
}
