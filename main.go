package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"tritchgo/config"

	"github.com/go-resty/resty/v2"
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

	log.Printf("Now: %v", time.Now())
	log.Printf("Token Exp: %v", tokenExpires)

	return token, nil
}

func main() {
	env := config.LoadEnv()
	client_id = env.CLIENT_ID
	client_secret = env.CLIENT_SECRET

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
			log.Printf("Stream %v", streams)
		}

		time.Sleep(30 * time.Second)
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
