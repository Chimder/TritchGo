package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"tritchgo/config"
	"tritchgo/internal/queries"
	"tritchgo/sqlc"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type twitchToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

var client_id, client_secret string

var (
	token        string
	mu           sync.Mutex
	client       = resty.New()
	tokenExpires time.Time
	isRefreshing bool
	// cond         = sync.NewCond(&mu)
)

func getToken() (string, time.Time, error) {
	tokenResp := &twitchToken{}
	resp, err := client.R().SetQueryParams(map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"grant_type":    "client_credentials",
	}).SetResult(tokenResp).Post("https://id.twitch.tv/oauth2/token")

	if err != nil || resp.StatusCode() != 200 || tokenResp.AccessToken == "" {
		log.Printf("Err fetch token %v", err)
		return "", time.Time{}, fmt.Errorf("Err fetch token")
	}

	safeExpirationTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-10 * time.Hour)

	return tokenResp.AccessToken, safeExpirationTime, nil
}

func getValidToken() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	log.Print("Start Fetch New Token")
	if !tokenExpires.IsZero() && time.Now().Before(tokenExpires) {
		log.Print("Token is still valid")
		return token, nil
	}

	newToken, expirationTime, err := getToken()
	if err != nil {
		return "", fmt.Errorf("Cant fetch new token: %v", err)
	}

	token = newToken
	tokenExpires = expirationTime

	return token, nil
}

func main() {
	context := context.Background()
	env := config.LoadEnv()
	client_id = env.CLIENT_ID
	client_secret = env.CLIENT_SECRET
	///////////////
	db, err := sqlc.DBConn(context)
	sqlc := queries.New(db)
	log.Printf("da", sqlc)
	///////////
	if err != nil {
		log.Fatalf("Fatal conn to db: %v", err)
	}

	for {
		_, err := getValidToken()
		topGames, err := GetTopGames()
		if err != nil {
			log.Println("Err fetch Top Games")
			return
		}
		for _, game := range topGames {
			streams, err := GetTopStream(game.ID)
			if err != nil {
				log.Println("Err fetch Top Games")
				return
			}
			log.Print("Range Steams")
			for _, stream := range streams {
				log.Printf("DB FETCHE STERAM: %v", stream.UserID)
				startedAt, err := time.Parse(time.RFC3339, stream.StartedAt)
				if err != nil {
					log.Println("Error parsing started_at:", err)
					// continue
				}
				now := time.Now().UTC()
				date := pgtype.Date{Time: now, Valid: true}

				airtimeMinutes := int(now.Sub(startedAt).Minutes())

				err = sqlc.InsertStreamStats(context, queries.InsertStreamStatsParams{
					StreamID:       stream.ID,
					UserID:         stream.UserID,
					GameID:         stream.GameID,
					Date:           date,
					Airtime:        pgtype.Int4{Int32: int32(airtimeMinutes), Valid: true},
					PeakViewers:    pgtype.Int4{Int32: int32(stream.ViewerCount), Valid: true},
					AverageViewers: pgtype.Int4{Int32: int32(stream.ViewerCount), Valid: true},
					HoursWatched:   pgtype.Int4{Int32: 0, Valid: true},
				})
				if err != nil {
					log.Println("Err fetch Top Games")
					return
				}
			}
		}

		time.Sleep(15 * time.Minute)
	}

}

func GetTopGames() ([]Game, error) {
	var authHeaders = map[string]string{
		"Client-ID":     client_id,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	topGames := &TopGamesResponse{}
	respTopGames, err := client.R().SetHeaders(authHeaders).SetQueryParam("first", strconv.Itoa(1)).SetResult(topGames).Get("https://api.twitch.tv/helix/games/top")

	if err != nil || respTopGames.StatusCode() != 200 {
		return nil, fmt.Errorf("Top Games fetch Err: %v", respTopGames.Error())
	}
	return topGames.Data, nil

}

func GetTopStream(gameId string) ([]Stream, error) {
	var authHeaders = map[string]string{
		"Client-ID":     client_id,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	var topStreamers = &StreamsResponse{}
	respTopStreamer, err := client.R().SetHeaders(authHeaders).SetQueryParams(map[string]string{
		"game_id": gameId,
		"first":   strconv.Itoa(3),
	}).SetResult(topStreamers).Get("https://api.twitch.tv/helix/streams")

	if err != nil || respTopStreamer.StatusCode() != 200 {
		return nil, fmt.Errorf("Top Stream fetch Err: %v", respTopStreamer.Error())
	}
	return topStreamers.Data, nil

}
