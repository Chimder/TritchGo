package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StreamStats struct {
	ID             uuid.UUID `json:"id" db:"id"`
	StreamID       string    `json:"stream_id" db:"stream_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	GameID         string    `json:"game_id" db:"game_id"`
	Date           time.Time `json:"date" db:"date"`
	Title          string    `json:"title" db:"title"`
	Airtime        int       `json:"airtime" db:"airtime"`
	PeakViewers    int       `json:"peak_viewers" db:"peak_viewers"`
	AverageViewers int       `json:"average_viewers" db:"average_viewers"`
	HoursWatched   int       `json:"hours_watched" db:"hours_watched"`
}

type StatsRepo interface {
	GetUserStatsById(ctx context.Context, id string) ([]StreamStats, error)
	GetStreamStatsById(ctx context.Context, id string) ([]StreamStats, error)
	GetAllStats() ([]StreamStats, error)
	BuildInsertQuery(stream *StreamArg, startedAt time.Time) (string, pgx.NamedArgs)
}

type statsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepo(db *pgxpool.Pool) StatsRepo {
	return &statsRepo{db: db}
}
func (r *statsRepo) GetAllStats() ([]StreamStats, error) {
	ctx := context.Background()
	query := `SELECT * FROM stream_stats`
	rows, err := r.db.Query(ctx, query)
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
func (r *statsRepo) BuildInsertQuery(stream *StreamArg, startedAt time.Time) (string, pgx.NamedArgs) {
	now := time.Now().UTC()
	airtimeMinutes := int(time.Since(startedAt).Minutes())

	query := `INSERT INTO stream_stats (
		stream_id, user_id, game_id, date, title, airtime, peak_viewers, average_viewers, hours_watched
	) VALUES (
		@stream_id, @user_id, @game_id, @date, @title, @airtime, @peak_viewers, @average_viewers, @hours_watched
	) ON CONFLICT (stream_id, date)
	DO UPDATE SET
		airtime = EXCLUDED.airtime,
		peak_viewers = GREATEST(stream_stats.peak_viewers, EXCLUDED.peak_viewers),
		average_viewers = ROUND((stream_stats.average_viewers + EXCLUDED.average_viewers) / 2.0),
		hours_watched = stream_stats.hours_watched + ROUND(EXCLUDED.average_viewers * (EXCLUDED.airtime / 60.0));`

	args := pgx.NamedArgs{
		"stream_id":       stream.ID,
		"user_id":         stream.UserID,
		"date":            now,
		"title":           stream.Title,
		"game_id":         stream.GameID,
		"peak_viewers":    stream.ViewerCount,
		"airtime":         airtimeMinutes,
		"average_viewers": stream.ViewerCount,
		"hours_watched":   0,
	}

	return query, args
}

type StreamArg struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	GameID      string `json:"game_id"`
	ViewerCount int    `json:"viewer_count"`
}
