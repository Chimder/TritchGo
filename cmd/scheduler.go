package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"sync"
	"time"
	"tritchgo/internal/handlers"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TwitchScheduler struct {
	ctx context.Context
	db  *pgxpool.Pool
	es  *elasticsearch.Client
}

func NewTwitchScheduler(ctx context.Context, db *pgxpool.Pool, es *elasticsearch.Client) *TwitchScheduler {
	return &TwitchScheduler{ctx: ctx, db: db, es: es}
}

func nextInterval(duration time.Duration) time.Time {
	now := time.Now()
	minutes := now.Minute()
	nextMinutes := (minutes/int(duration.Minutes()) + 1) * int(duration.Minutes())
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextMinutes, 0, 0, now.Location())
}

func (ts *TwitchScheduler) StartFetchLoop(twitchHandle *handlers.TwitchHandle) {
	for {
		interval := 1 * time.Minute
		nextTick := nextInterval(interval)

		log.Printf("Next TICK: %v", nextTick)
		time.Sleep(time.Until(nextTick))

		err := ts.fetchAndStoreTopGames(twitchHandle)
		if err != nil {
			slog.Error("FetchAndStoreTopGames:", "Error", err)
		}

		time.Sleep(3 * time.Minute)
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

	for streams := range gameChan {
		for _, stream := range streams {
			insertWg.Add(1)
			go func(stream handlers.Stream) {
				defer insertWg.Done()

				startedAt, err := time.Parse(time.RFC3339, stream.StartedAt)
				if err != nil {
					log.Println("Error parsing started_at:", err)
					return
				}

				batchMutex.Lock()
				ts.insertStreamStats(&stream, streamBatchInsert, startedAt)
				batchMutex.Unlock()

				err = ts.indexStreamToElastic(&stream)
				if err != nil {
					log.Printf("Failed to index to Elastic: %v", err)
				}
			}(stream)
		}
	}

	insertWg.Wait()

	res := ts.db.SendBatch(ts.ctx, streamBatchInsert)
	defer res.Close()

	for i := 0; i < streamBatchInsert.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Printf("Postgres error: %v", pgErr)
			}
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	slog.Info("timeTaken", "time", time.Since(startTime))
	return nil
}

func (ts *TwitchScheduler) insertStreamStats(stream *handlers.Stream, streamBatchInsert *pgx.Batch, startedAt time.Time) {
	now := time.Now().UTC()
	airtimeDuration := now.Sub(startedAt)
	airtimeMinutes := int(math.Round(airtimeDuration.Minutes()))

	stringQuery := `INSERT INTO stream_stats (
    stream_id, user_id, game_id, date, title, airtime, peak_viewers, average_viewers, hours_watched
) VALUES (
  @stream_id, @user_id, @game_id, @date, @title, @airtime, @peak_viewers, @average_viewers, @hours_watched
) ON CONFLICT (stream_id, date)
DO UPDATE SET
    airtime = EXCLUDED.airtime,
    peak_viewers = GREATEST(stream_stats.peak_viewers, EXCLUDED.peak_viewers),
    average_viewers = ROUND(stream_stats.average_viewers + EXCLUDED.average_viewers) / 2,
    hours_watched = stream_stats.hours_watched + ROUND(EXCLUDED.average_viewers * (EXCLUDED.airtime / 60.0)); `
	args := pgx.NamedArgs{
		"stream_id":       stream.ID,
		"user_id":         stream.UserID,
		"date":            now,
		"title":           stream.Title,
		"game_id":         stream.GameID,
		"peak_viewers":    stream.ViewerCount,
		"airtime":         airtimeMinutes,
		"average_viewers": stream.ViewerCount,
		"hours_watched":   0,
	}

	streamBatchInsert.Queue(stringQuery, args)
}

func (ts *TwitchScheduler) indexStreamToElastic(stream *handlers.Stream) error {
	doc := map[string]interface{}{
		"user_id": stream.UserID,
		"title":   stream.Title,
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := bytes.NewReader(docBytes)

	res, err := ts.es.Index(
		"stream_stats",
		req,
		ts.es.Index.WithDocumentID(stream.ID),
		ts.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch indexing error: %s", res.String())
	}

	return nil
}
