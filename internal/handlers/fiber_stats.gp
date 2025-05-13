package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"tritchgo/internal/repository"
	"tritchgo/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type StatsHandler struct {
	pgdb *pgxpool.Pool
	repo *repository.Repository
	rdb  *redis.Client
	// kafkaWriter *kafkaW.KafkaWriters
}

// func NewStatsHandler(repo *repository.Repository, db *pgxpool.Pool, rdb *redis.Client, kafkaWriter *kafkaW.KafkaWriters) *StatsHandler {
func NewStatsHandler(repo *repository.Repository, db *pgxpool.Pool, rdb *redis.Client) *StatsHandler {
	return &StatsHandler{
		repo: repo,
		pgdb: db,
		rdb:  rdb,
		// kafkaWriter: kafkaWriter,
	}
}

func (st *StatsHandler) Test(w http.ResponseWriter, r *http.Request) {
	newmap := map[string]string{"test1": "lox", "test2": "lox2"}
	cache, err := st.rdb.HGetAll(r.Context(), "test").Result()
	if err != nil {
		log.Print("Redis error:", err)
	}

	if len(cache) > 0 {
		log.Print("from cache")
		utils.WriteJSON(w, 200, cache)
		return
	}

	log.Print("start set")
	err = st.rdb.HSet(r.Context(), "test", newmap).Err()
	if err != nil {
		utils.WriteError(w, 500, "Err fetch user stats")
		return
	}
	st.rdb.Expire(r.Context(), "test", 30*time.Second)
	log.Print("form map")
	utils.WriteJSON(w, 200, newmap)
}

func (st *StatsHandler) GetUserStatsById(c *fiber.Ctx) error {
	userId := c.Query("user_id")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id required",
		})
	}

	cacheData, err := st.rdb.Get(c.Context(), userId).Result()
	if err == nil {
		return c.Status(fiber.StatusOK).Type("json").Send([]byte(cacheData))
	}

	stats, err := st.repo.Stats.GetUserStatsById(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(stats) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "err marshal data",
		})
	}

	// err = st.kafkaWriter.UserStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(userId), Value: data})
	// if err != nil {
	// 	log.Printf("Kafka Err %v", err)
	// }

	if err := st.rdb.Set(c.Context(), userId, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	return c.Status(fiber.StatusOK).Type("json").Send(data)
}

func (st *StatsHandler) GetStreamStatsById(c *fiber.Ctx) error {
	stream_id := c.Query("stream_id")
	if stream_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "stream_id required",
		})
	}

	cacheData, err := st.rdb.Get(c.Context(), stream_id).Result()
	if err == nil {
		return c.Status(fiber.StatusOK).Type("json").Send([]byte(cacheData))
	}

	stats, err := st.repo.Stats.GetStreamStatsById(c.Context(), stream_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "err get stream from db",
		})
	}
	if len(stats) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "stream not found",
		})
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "err marshal data",
		})
	}

	// err = st.kafkaWriter.StreamStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(stream_id), Value: data})
	// if err != nil {
	// 	log.Printf("Kafka Err %v", err)
	// }

	if err := st.rdb.Set(c.Context(), stream_id, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	return c.Status(fiber.StatusOK).Type("json").Send(data)
}
