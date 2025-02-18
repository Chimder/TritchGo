package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
	"tritchgo/internal/handlers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TwitchSheduler struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewTwitchSheduler(ctx context.Context, db *pgxpool.Pool) *TwitchSheduler {
	return &TwitchSheduler{ctx: ctx, db: db}
}

func nextInterval(duration time.Duration) time.Time {
	now := time.Now()
	minutes := now.Minute()
	nextMinutes := (minutes/int(duration.Minutes()) + 1) * int(duration.Minutes())
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextMinutes, 0, 0, now.Location())
}

func (ts *TwitchSheduler) StartFetchLoop(twitchHandle *handlers.TwitchHandle) {
	for {
		interval := 3 * time.Minute
		nextTick := nextInterval(interval).Add(-1 * time.Minute)

		log.Printf("Next TICK: %v", nextTick)
		time.Sleep(time.Until(nextTick))

		err := ts.fetchAndStoreTopGames(twitchHandle)
		if err != nil {
			log.Println("Error in fetchAndStoreTopGames:", err)
		}

		time.Sleep(2 * time.Minute)
	}
}

func (ts *TwitchSheduler) fetchAndStoreTopGames(twitchHandle *handlers.TwitchHandle) error {
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

	gameChan := make(chan []handlers.Stream, 100)
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
		log.Print("FirstWAITE CLOSe")
		close(gameChan)
	}()

	var insertWg sync.WaitGroup
	streamBatchInsert := &pgx.Batch{}
	for streams := range gameChan {
		for _, stream := range streams {
			insertWg.Add(1)
			go func(stream handlers.Stream) {
				defer insertWg.Done()

				startedAt, err := time.Parse(time.RFC3339, stream.StartedAt)
				if err != nil {
					log.Println("Error parsing started_at:", err)
				}
				ts.insertStreamStats(&stream, streamBatchInsert, startedAt)
			}(stream)
		}
	}
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

	// log.Printf("BATCHRESl %v", res)
	insertWg.Wait()
	log.Print("SECCClllllWAITE CLOSE")
	elapsed := time.Since(startTime)
	log.Printf("time taken: %v", elapsed)

	return nil
}

func (ts *TwitchSheduler) insertStreamStats(stream *handlers.Stream, streamBatchInsert *pgx.Batch, startedAt time.Time) {
	now := time.Now().UTC()
	airtimeDuration := now.Sub(startedAt)
	airtimeMinutes := int(math.Round(airtimeDuration.Minutes()))

	stringQuery := `INSERT INTO stream_stats (
    stream_id, user_id, game_id, date, airtime, peak_viewers, average_viewers, hours_watched
) VALUES (
  @stream_id, @user_id, @game_id, @date, @airtime, @peak_viewers, @average_viewers, @hours_watched
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
		"game_id":         stream.GameID,
		"peak_viewers":    stream.ViewerCount,
		"airtime":         airtimeMinutes,
		"average_viewers": stream.ViewerCount,
		"hours_watched":   0,
	}

	streamBatchInsert.Queue(stringQuery, args)
	// _, err := ts.db.Exec(ts.ctx, stringQuery, args)
	// if err != nil {
	// 	log.Printf("Err Set Top Games to db %v", err)
	// }
}
