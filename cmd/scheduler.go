package main

import (
	"context"
	"log"
	"math"
	"time"
	"tritchgo/internal/handlers"

	"github.com/jackc/pgx/v5"
)

func nextInterval(duration time.Duration) time.Time {
	now := time.Now()
	minutes := now.Minute()
	nextMinutes := (minutes/int(duration.Minutes()) + 1) * int(duration.Minutes())
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextMinutes, 0, 0, now.Location())
}

func StartFetchLoop(ctx context.Context, twitchHandle *handlers.TwitchHandle, db *pgx.Conn) {
	for {
		interval := 10 * time.Minute
		nextTick := nextInterval(interval).Add(-1 * time.Minute)

		log.Printf("Next TICK: %v", nextTick)
		time.Sleep(time.Until(nextTick))

		err := fetchAndStoreTopGames(ctx, twitchHandle, db)
		if err != nil {
			log.Println("Error in fetchAndStoreTopGames:", err)
		}

		time.Sleep(2 * time.Minute)
	}
}

func fetchAndStoreTopGames(ctx context.Context, twitchHandle *handlers.TwitchHandle, db *pgx.Conn) error {
	_, err := twitchHandle.GetValidToken()
	if err != nil {
		return err
	}

	log.Print("Start Fetch")
	topGames, err := twitchHandle.GetTopGames()
	if err != nil {
		return err
	}

	for _, game := range topGames {
		streams, err := twitchHandle.GetTopStream(game.ID)
		if err != nil {
			return err
		}
		for _, stream := range streams {
			startedAt, err := time.Parse(time.RFC3339, stream.StartedAt)
			if err != nil {
				log.Println("Error parsing started_at:", err)
				continue
			}

			now := time.Now().UTC()
			airtimeDuration := now.Sub(startedAt)
			airtimeMinutes := int(math.Round(airtimeDuration.Minutes()))

			stringQuery := `INSERT INTO stream_stats (
    stream_id, user_id, game_id, date, airtime, peak_viewers, average_viewers, hours_watched
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (stream_id, date)
DO UPDATE SET
    airtime = EXCLUDED.airtime,
    peak_viewers = GREATEST(stream_stats.peak_viewers, EXCLUDED.peak_viewers),
    average_viewers = ROUND(stream_stats.average_viewers + EXCLUDED.average_viewers) / 2,
    hours_watched = stream_stats.hours_watched + ROUND(EXCLUDED.average_viewers * (EXCLUDED.airtime / 60.0)); `
			rows, err := db.Exec(ctx, stringQuery, stream.ID,
				stream.UserID,
				stream.GameID,
				now,
				airtimeMinutes,
				stream.ViewerCount,
				stream.ViewerCount,
				0,
			)
			rowsAffected := rows.RowsAffected()
			log.Printf("success: %d", rowsAffected)
			if err != nil {
				log.Printf("Err Set Top Games to db %v", err)
				return err
			}
		}
	}

	return nil
}
