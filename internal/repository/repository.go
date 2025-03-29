package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Stats StatsRepo
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Stats: NewStatsRepo(db),
	}
}
