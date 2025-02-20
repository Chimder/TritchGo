package handlers

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsHandler struct {
	db *pgxpool.Pool
}

func NewStatsHandler(db *pgxpool.Pool) *StatsHandler {
	return &StatsHandler{
		db: db,
	}
}

type StreamStats struct {
	ID             int       `json:"id" db:"id"`
	StreamID       string    `json:"stream_id" db:"stream_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	GameID         string    `json:"game_id" db:"game_id"`
	Date           time.Time `json:"date" db:"date"`
	Airtime        int       `json:"airtime" db:"airtime"`
	PeakViewers    int       `json:"peak_viewers" db:"peak_viewers"`
	AverageViewers int       `json:"average_viewers" db:"average_viewers"`
	HoursWatched   int       `json:"hours_watched" db:"hours_watched"`
}

func (st *StatsHandler) GetUserStatsById(ctx context.Context, id string) ([]StreamStats, error) {
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

func (st *StatsHandler) GetStreamStatsById(ctx context.Context, id string) ([]StreamStats, error) {
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
