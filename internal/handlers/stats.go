package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"tritchgo/internal/store"
	"tritchgo/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type StatsHandler struct {
	pgdb  *pgxpool.Pool
	store *store.Storage
	rdb   *redis.Client
}

func NewStatsHandler(db *pgxpool.Pool, rdb *redis.Client) *StatsHandler {
	store := store.NewStorage(db)
	return &StatsHandler{
		store: &store,
		pgdb:  db,
		rdb:   rdb,
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

	cacheData, err := st.rdb.Get(r.Context(), userId).Result()
	if err == nil {
		utils.WriteJSONRedis(w, 200, []byte(cacheData))
		return
	}

	stats, err := st.store.Stats.GetUserStatsById(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, 500, "Err fetch user stats")
		return
	}
	data, err := json.Marshal(stats)
	log.Printf("marshal after fetch %v ", data)
	if err != nil {
		log.Printf("Error json stats: %v", err)
		utils.WriteError(w, 500, "Failed to marsh data")
		return
	}
	if err := st.rdb.Set(r.Context(), userId, data, 30*time.Second).Err(); err != nil {
		log.Printf("Err cache user stats: %v", err)
	}
	utils.WriteJSONRedis(w, 200, data)
}

func (st *StatsHandler) GetStreamStatsById(w http.ResponseWriter, r *http.Request) {
	streamId := r.URL.Query().Get("stream_id")
	redisdata, err := st.rdb.Get(r.Context(), streamId).Result()
	if err == nil {
		utils.WriteJSONRedis(w, 200, []byte(redisdata))
		return
	}

	stats, err := st.store.Stats.GetStreamStatsById(r.Context(), streamId)
	if err != nil {
		log.Printf("Err encode user stats  %v", err)
		return
	}
	data, err := json.Marshal(stats)
	log.Printf("marshal after fetch %v ", data)

	if err != nil {
		log.Printf("Error marshal stats: %v", err)
		utils.WriteError(w, 500, "Failed to marsh data")
		return
	}
	if err = st.rdb.Set(r.Context(), streamId, data, 30*time.Second).Err(); err != nil {
		log.Printf("Err cache user stats: %v", err)
	}
	utils.WriteJSONRedis(w, 200, data)
}
