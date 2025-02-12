package main

import (
	"context"
	"log"
	"math"
	"time"
	"tritchgo/internal/handlers"
	"tritchgo/internal/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

func nextInterval(duration time.Duration) time.Time {
	now := time.Now()
	minutes := now.Minute()
	nextMinutes := (minutes/int(duration.Minutes()) + 1) * int(duration.Minutes())
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), nextMinutes, 0, 0, now.Location())
}

func StartFetchLoop(ctx context.Context, twitchHandle *handlers.TwitchHandle, sqlc *queries.Queries) {
	for {
		interval := 15 * time.Minute
		nextTick := nextInterval(interval).Add(-1 * time.Minute)

		log.Printf("Next TICK: %v", nextTick)
		time.Sleep(time.Until(nextTick))

		err := fetchAndStoreTopGames(ctx, twitchHandle, sqlc)
		if err != nil {
			log.Println("Error in fetchAndStoreTopGames:", err)
		}

		time.Sleep(2 * time.Minute)
	}
}

func fetchAndStoreTopGames(ctx context.Context, twitchHandle *handlers.TwitchHandle, sqlc *queries.Queries) error {
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

			err = sqlc.InsertStreamStats(ctx, queries.InsertStreamStatsParams{
				StreamID:       stream.ID,
				UserID:         stream.UserID,
				GameID:         stream.GameID,
				Date:           pgtype.Date{Time: now, Valid: true},
				Airtime:        pgtype.Int4{Int32: int32(airtimeMinutes), Valid: true},
				PeakViewers:    pgtype.Int4{Int32: int32(stream.ViewerCount), Valid: true},
				AverageViewers: pgtype.Int4{Int32: int32(stream.ViewerCount), Valid: true},
				HoursWatched:   pgtype.Int4{Int32: 0, Valid: true},
			})
			if err != nil {
				log.Printf("Err Set Top Games to db %v", err)
				return err
			}
		}
	}

	return nil
}
