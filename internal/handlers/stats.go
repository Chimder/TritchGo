package handlers

import (
	"log"
	"net/http"
	"time"
	"tritchgo/internal/store"
	"tritchgo/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsHandler struct {
	db    *pgxpool.Pool
	store *store.Storage
}

func NewStatsHandler(db *pgxpool.Pool) *StatsHandler {
	store := store.NewStorage(db)
	return &StatsHandler{
		store: &store,
		db:    db,
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

func (st *StatsHandler) GetUserStatsById(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	stats, err := st.store.Stats.GetUserStatsById(r.Context(), userId)
	if err != nil {
		log.Printf("Err fetch user stats  %v", err)
		return
	}

	utils.WriteJSON(w, 200, stats)
}

func (st *StatsHandler) GetStreamStatsById(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream_id")

	stats, err := st.store.Stats.GetStreamStatsById(r.Context(), streamId)
	if err != nil {
		log.Printf("Err encode user stats  %v", err)
		return
	}

	utils.WriteJSON(w, 200, stats)

}
