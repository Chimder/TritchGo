package routers

import (
	"net/http"
	"tritchgo/internal/handlers"
	kafkaWriter "tritchgo/internal/kafka"
	"tritchgo/internal/repository"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(repo *repository.Repository, pgdb *pgxpool.Pool, rdb *redis.Client, kafkaWriter *kafkaWriter.KafkaWriters, els *elasticsearch.Client) *chi.Mux {
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

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	statsHandle := handlers.NewStatsHandler(repo, pgdb, rdb, kafkaWriter)

	r.Get("/user/stats", statsHandle.GetUserStatsById)
	r.Get("/stream/stats", statsHandle.GetStreamStatsById)
	r.Get("/redis", statsHandle.Test)

	return r
}
