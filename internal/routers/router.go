package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"tritchgo/internal/queries"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(sqlc *queries.Queries) *chi.Mux {
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

	r.Route("/user", func(r chi.Router) {
		r.Mount("/", UserRouter(sqlc))
	})
	// r.Mount("/", GameRoutes())
	r.Get("/stream/stats", func(w http.ResponseWriter, r *http.Request) {
		streamId := r.URL.Query().Get("stream_id")
		streamStats, err := sqlc.GetStatsByStreamId(r.Context(), streamId)
		if err != nil {
			log.Printf("Err fetch stream stats: %v", err)
			return
		}
		err = json.NewEncoder(w).Encode(streamStats)
		if err != nil {
			log.Printf("Err encode user stats: %v", err)
			return
		}
	})

	return r
}
