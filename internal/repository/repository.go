package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Stats *StatsRepo
}

func NewRepository(db *pgxpool.Pool) *Repository {
	// NewStatsRepo(db)
	return &Repository{
		Stats: NewStatsRepo(db),
	}
}
