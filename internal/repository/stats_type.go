package repository

import (
	"time"

	"github.com/google/uuid"
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
