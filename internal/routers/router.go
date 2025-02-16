package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"tritchgo/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewRouter(db *pgxpool.Pool) *chi.Mux {
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

	statsHandle := handlers.NewStatsHandler(db)
	r.Get("/user/stats", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")

		stats, err := statsHandle.GetUserStatsById(r.Context(), userId)
		if err != nil {
			log.Printf("Err fetch user stats  %v", err)
			return
		}

		err = json.NewEncoder(w).Encode(stats)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}
	})

	r.Get("/stream/stats", func(w http.ResponseWriter, r *http.Request) {
		streamId := r.URL.Query().Get("stream_id")

		stats, err := statsHandle.GetStreamStatsById(r.Context(), streamId)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}

		err = json.NewEncoder(w).Encode(stats)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}
	})

	return r
}
