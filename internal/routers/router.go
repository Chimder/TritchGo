package routers

import (
	"net/http"
	"tritchgo/internal/handlers"
	kafkaWriter "tritchgo/internal/kafka"
	"tritchgo/internal/repository"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(repo *repository.Repository, pgdb *pgxpool.Pool, rdb *redis.Client, kafkaWriter *kafkaWriter.KafkaWriters, els *elasticsearch.Client) *gin.Engine {
	r := gin.Default()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running")
	})

	statsHandle := handlers.NewStatsHandler(repo, pgdb, rdb, kafkaWriter)

	api := r.Group("/")
	{
		api.GET("/user/stats", statsHandle.GetUserStatsById)
		api.GET("/stream/stats", statsHandle.GetStreamStatsById)
	}

	return r
}
