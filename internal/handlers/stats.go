package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	kafkaW "tritchgo/internal/kafka"
	"tritchgo/internal/repository"
	"tritchgo/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type StatsHandler struct {
	pgdb        *pgxpool.Pool
	repo        *repository.Repository
	rdb         *redis.Client
	kafkaWriter *kafkaW.KafkaWriters
}

func NewStatsHandler(repo *repository.Repository, db *pgxpool.Pool, rdb *redis.Client, kafkaWriter *kafkaW.KafkaWriters) *StatsHandler {
	return &StatsHandler{
		repo:        repo,
		pgdb:        db,
		rdb:         rdb,
		kafkaWriter: kafkaWriter,
	}
}

func (st *StatsHandler) Test(c *gin.Context) {
	newmap := map[string]string{"test1": "lox", "test2": "lox2"}
	cache, err := st.rdb.HGetAll(c.Request.Context(), "test").Result()
	if err != nil {
		log.Print("Redis error:", err)
	}

	if len(cache) > 0 {
		log.Print("from cache")
		utils.WriteJSON(c, 200, cache)
		return
	}

	log.Print("start set")
	err = st.rdb.HSet(c.Request.Context(), "test", newmap).Err()
	if err != nil {
		utils.WriteError(c, 500, "Err fetch user stats")
		return
	}
	st.rdb.Expire(c.Request.Context(), "test", 30*time.Second)
	log.Print("form map")
	utils.WriteJSON(c, 200, newmap)
}

func (st *StatsHandler) GetUserStatsById(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		utils.WriteError(c, 400, "user_id required")
		return
	}

	cacheData, err := st.rdb.Get(c.Request.Context(), userId).Result()
	if err == nil {
		utils.WriteJSONRedis(c, 200, []byte(cacheData))
		return
	}

	stats, err := st.repo.Stats.GetUserStatsById(c.Request.Context(), userId)
	if err != nil {
		utils.WriteError(c, 400, err.Error())
		return
	}
	if len(stats) == 0 {
		utils.WriteError(c, http.StatusNotFound, "user not found")
		return
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		utils.WriteError(c, 400, "err marshal data")
		return
	}

	// err = st.kafkaWriter.UserStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(userId), Value: data})
	// if err != nil {
	// 	log.Printf("Kafka Err %v", err)
	// }

	if err := st.rdb.Set(c.Request.Context(), userId, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	utils.WriteJSONRedis(c, 200, data)
}

func (st *StatsHandler) GetStreamStatsById(c *gin.Context) {
	stream_id := c.Query("stream_id")
	if stream_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "stream_id required"})
		return
	}

	cacheData, err := st.rdb.Get(c.Request.Context(), stream_id).Result()
	if err == nil {
		utils.WriteJSONRedis(c, 200, []byte(cacheData))
		return
	}

	stats, err := st.repo.Stats.GetStreamStatsById(c.Request.Context(), stream_id)
	if err != nil {
		utils.WriteError(c, 400, "err get stream from db")
		return
	}
	if len(stats) == 0 {
		utils.WriteError(c, http.StatusNotFound, "stream not found")
		return
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		utils.WriteError(c, 400, "err marshal data")
		return
	}

	// err = st.kafkaWriter.StreamStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(stream_id), Value: data})
	// if err != nil {
	// 	log.Printf("Kafka Err %v", err)
	// }

	if err := st.rdb.Set(c.Request.Context(), stream_id, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	utils.WriteJSONRedis(c, 200, data)
}
