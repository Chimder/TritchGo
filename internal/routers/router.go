package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

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

func NewRouter(db *pgx.Conn) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	r.Get("/user/stats", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")

		query := `SELECT * FROM stream_stats WHERE user_id = $1`
		rows, err := db.Query(r.Context(), query, userId)
		if err != nil {
			log.Printf("Err fetch user stats  %v", err)
			return
		}
		defer rows.Close()
		var stats []StreamStats
		for rows.Next() {
			var s StreamStats
			err := rows.Scan(&s.ID, &s.StreamID, &s.UserID, &s.GameID, &s.Date, &s.Airtime, &s.PeakViewers, &s.AverageViewers, &s.HoursWatched)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			stats = append(stats, s)
		}
		if err := rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
		}

		err = json.NewEncoder(w).Encode(stats)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}
	})

	r.Get("/stream/stats", func(w http.ResponseWriter, r *http.Request) {
		streamId := r.URL.Query().Get("stream_id")

		query := `SELECT * FROM stream_stats WHERE stream_id = $1`
		var s StreamStats
		err := db.QueryRow(r.Context(), query, streamId).Scan(
			&s.ID, &s.StreamID, &s.UserID, &s.GameID, &s.Date, &s.Airtime, &s.PeakViewers, &s.AverageViewers, &s.HoursWatched)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Stream not found", http.StatusNotFound)
				return
			}
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(s)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}
	})

	return r
}
