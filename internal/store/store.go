package store

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Stats *StatsStore
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Stats: &StatsStore{db},
	}
}
