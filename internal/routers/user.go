package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"tritchgo/internal/queries"

	"github.com/go-chi/chi"
)

func UserRouter(sqlc *queries.Queries) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Route"))
	})

	r.Get("/all/stats", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")

		userStats, err := sqlc.GetStatsByUserId(r.Context(), userId)
		if err != nil {
			log.Printf("Err fetch user stats  %v", err)
			return
		}

		err = json.NewEncoder(w).Encode(userStats)
		if err != nil {
			log.Printf("Err encode user stats  %v", err)
			return
		}

	})

	return r
}
