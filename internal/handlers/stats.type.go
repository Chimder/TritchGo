package handlers

import (
	"time"
	"tritchgo/internal/repository"

	"github.com/google/uuid"
)

type StreamStatsResp struct {
	ID             uuid.UUID `json:"id"`
	StreamID       string    `json:"stream_id"`
	UserID         string    `json:"user_id"`
	GameID         string    `json:"game_id"`
	Date           time.Time `json:"date"`
	Airtime        int       `json:"airtime"`
	PeakViewers    int       `json:"peak_viewers"`
	AverageViewers int       `json:"average_viewers"`
	HoursWatched   int       `json:"hours_watched"`
}

func StreamStatsRespFromDB(s repository.StreamStats) StreamStatsResp {
	return StreamStatsResp{
		ID:             s.ID,
		StreamID:       s.StreamID,
		UserID:         s.UserID,
		GameID:         s.GameID,
		Date:           s.Date,
		Airtime:        s.Airtime,
		PeakViewers:    s.PeakViewers,
		AverageViewers: s.AverageViewers,
		HoursWatched:   s.HoursWatched,
	}
}

func StreamStatsResponseListFromDB(stats []repository.StreamStats) []StreamStatsResp {
	responses := make([]StreamStatsResp, len(stats))
	for i, s := range stats {
		responses[i] = StreamStatsRespFromDB(s)
	}
	return responses
}
