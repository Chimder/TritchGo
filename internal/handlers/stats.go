package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	kafkaW "tritchgo/internal/kafka"
	"tritchgo/internal/repository"
	"tritchgo/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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

func (st *StatsHandler) GetUserStatsById(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		utils.WriteError(w, 400, "user_id required")
		return
	}

	cacheData, err := st.rdb.Get(r.Context(), userId).Result()
	if err == nil {
		utils.WriteJSONRedis(w, 200, []byte(cacheData))
		return
	}

	stats, err := st.repo.Stats.GetUserStatsById(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, 400, err.Error())
		return
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		utils.WriteError(w, 400, "err marshal data")
		return
	}

	err = st.kafkaWriter.UserStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(userId), Value: data})
	if err != nil {
		log.Printf("Kafka Err %v", err)
	}

	if err := st.rdb.Set(r.Context(), userId, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	utils.WriteJSONRedis(w, 200, data)
}

func (st *StatsHandler) GetStreamStatsById(w http.ResponseWriter, r *http.Request) {
	stream_id := r.URL.Query().Get("stream_id")
	if stream_id == "" {
		utils.WriteError(w, 400, "user_id required")
		return
	}

	cacheData, err := st.rdb.Get(r.Context(), stream_id).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cacheData))
		return
	}

	stats, err := st.repo.Stats.GetStreamStatsById(r.Context(), stream_id)
	if err != nil {
		utils.WriteError(w, 400, "err get stream from db")
		return
	}

	resp := StreamStatsResponseListFromDB(stats)
	data, err := json.Marshal(resp)
	if err != nil {
		utils.WriteError(w, 400, "err marshal data")
		return
	}

	err = st.kafkaWriter.StreamStatsWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(stream_id), Value: data})
	if err != nil {
		log.Printf("Kafka Err %v", err)
	}

	if err := st.rdb.Set(r.Context(), stream_id, data, 2*time.Minute).Err(); err != nil {
		log.Printf("Err set user_id %v", err)
	}

	utils.WriteJSONRedis(w, 200, data)
}
