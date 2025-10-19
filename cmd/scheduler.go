package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"
	"tritchgo/internal/handlers"
	"tritchgo/internal/nats"
	"tritchgo/internal/repository"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TwitchScheduler struct {
	ctx  context.Context
	db   *pgxpool.Pool
	es   *elasticsearch.Client
	nc   *nats.NatsProducer
	repo *repository.Repository
}

func NewTwitchScheduler(ctx context.Context, db *pgxpool.Pool, es *elasticsearch.Client, nc *nats.NatsProducer, repo *repository.Repository) *TwitchScheduler {
	return &TwitchScheduler{ctx: ctx, db: db, es: es, nc: nc, repo: repo}
}

func NextInterval(interval int) time.Duration {
	return NextIntervalAt(time.Now(), interval)
}

func NextIntervalAt(now time.Time, interval int) time.Duration {
	currentMin := now.Minute()
	remaind := currentMin % interval
	nextMin := interval - remaind
	return time.Duration(nextMin) * time.Minute
}

func (ts *TwitchScheduler) StartFetchLoop(twitchHandle *handlers.TwitchHandle) {
	for {
		waitTick := NextInterval(1)

		log.Printf("Wait tick: %v", waitTick)
		time.Sleep(waitTick)

		err := ts.fetchAndStoreTopGames(twitchHandle)
		if err != nil {
			slog.Error("FetchAndStoreTopGames:", "Error", err)
		}

	}
}

func (ts *TwitchScheduler) fetchAndStoreTopGames(twitchHandle *handlers.TwitchHandle) error {
	_, err := twitchHandle.GetValidToken()
	if err != nil {
		return err
	}

	log.Print("Start Fetch")
	startTime := time.Now()

	topGames, err := twitchHandle.GetTopGames()
	if err != nil {
		return err
	}

	gameChan := make(chan []handlers.Stream, 55)
	var wg sync.WaitGroup
	for _, game := range topGames {
		wg.Add(1)
		go func(gameId string) {
			defer wg.Done()

			streams, err := twitchHandle.GetTopStream(gameId)
			if err != nil {
				log.Println("Error getTopStream:", err)
				return
			}

			gameChan <- streams
		}(game.ID)
	}

	go func() {
		wg.Wait()
		close(gameChan)
	}()

	var insertWg sync.WaitGroup
	var batchMutex sync.Mutex

	streamBatchInsert := &pgx.Batch{}

	var elasticBuf bytes.Buffer
	var elasticMutex sync.Mutex

	type NatsStreamData struct {
		StreamID string `json:"stream_id"`
		UserID   string `json:"user_id"`
	}
	for streams := range gameChan {
		for _, stream := range streams {
			insertWg.Add(1)

			go func(stream *handlers.Stream) {
				if stream == nil {
					return
				}
				defer insertWg.Done()

				// nats
				ncdata, err := json.Marshal(NatsStreamData{StreamID: stream.ID, UserID: stream.UserID})
				if err != nil {
					log.Println("Error marshaling natsData: %w", err)
				}

				if err := ts.nc.Publish(ts.ctx, "tritch.stats", ncdata); err != nil {
					log.Println("Error publishing to NATS:", err)
				}
				//

				startedAt, err := time.Parse(time.RFC3339, stream.StartedAt)
				if err != nil {
					log.Println("Error parsing started_at: %w", err)
					return
				}

				queryStr, nameArg := ts.repo.Stats.BuildInsertQuery(&repository.StreamArg{ID: stream.ID,
					UserID:      stream.UserID,
					Title:       stream.Title,
					GameID:      stream.GameID,
					ViewerCount: stream.ViewerCount},
					startedAt)

				batchMutex.Lock()
				streamBatchInsert.Queue(queryStr, nameArg)
				batchMutex.Unlock()

				elasticMutex.Lock()
				ts.indexStreamToElastic(stream, &elasticBuf)
				elasticMutex.Unlock()
			}(&stream)
		}
	}

	insertWg.Wait()

	elasticRes, err := ts.es.Bulk(
		strings.NewReader(elasticBuf.String()),
		ts.es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		log.Printf("Error sending bulk request: %s", err)
	}
	defer elasticRes.Body.Close()
	log.Printf("Bulk response: %s", elasticRes.Status())

	if streamBatchInsert.Len() == 0 {
		return fmt.Errorf("No queries queued in batch %w", err)
	}

	res := ts.db.SendBatch(ts.ctx, streamBatchInsert)
	defer res.Close()

	for i := 0; i < streamBatchInsert.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return fmt.Errorf("postgres batch exec: %w", err)
		}
	}

	slog.Info("timeTaken", "time", time.Since(startTime))
	return nil
}

func (ts *TwitchScheduler) indexStreamToElastic(stream *handlers.Stream, buf *bytes.Buffer) {
	meta := map[string]interface{}{
		"index": map[string]interface{}{
			"_index": "stream_stats",
			"_id":    stream.ID,
		},
	}
	doc := map[string]interface{}{
		"user_id": stream.UserID,
		"title":   stream.Title,
	}
	if err := json.NewEncoder(buf).Encode(meta); err != nil {
		log.Printf("Error encoding meta: %s", err)
	}

	if err := json.NewEncoder(buf).Encode(doc); err != nil {
		log.Printf("Error encoding document: %s", err)
	}
}
