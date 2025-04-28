package handlers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"tritchgo/config"

	"github.com/go-resty/resty/v2"
)

type TwitchHandle struct {
	client_id     string
	client_secret string
	mu            sync.Mutex
}

var (
	client       = resty.New()
	tokenExpires time.Time
	token        string
)

func NewTwitchHandle() *TwitchHandle {
	// env := config.LoadEnv()
	return &TwitchHandle{
		client_id:     config.LoadEnv().ClientID,
		client_secret: config.LoadEnv().ClientSecret,
	}
}

func (t *TwitchHandle) GetTopGames() ([]Game, error) {
	var authHeaders = map[string]string{
		"Client-ID":     t.client_id,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	topGames := &TopGamesResponse{}
	respTopGames, err := client.R().SetHeaders(authHeaders).SetQueryParam("first", strconv.Itoa(50)).SetResult(topGames).Get("https://api.twitch.tv/helix/games/top")

	if err != nil || respTopGames.StatusCode() != 200 {
		log.Printf("Unexpected status code: %d, response: %s", respTopGames.StatusCode(), respTopGames.String())
		return nil, fmt.Errorf("top Games fetch Err: status code %d", respTopGames.StatusCode())
	}
	if len(topGames.Data) == 0 {
		log.Println("No games returned from API")
		return nil, errors.New("no top games found")
	}

	return topGames.Data, nil

}

func (t *TwitchHandle) GetTopStream(gameId string) ([]Stream, error) {
	var authHeaders = map[string]string{
		"Client-ID":     t.client_id,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	var topStreamers = &StreamsResponse{}
	respTopStreamer, err := client.R().SetHeaders(authHeaders).SetQueryParams(map[string]string{
		"game_id": gameId,
		"first":   strconv.Itoa(100),
	}).SetResult(topStreamers).Get("https://api.twitch.tv/helix/streams")

	if err != nil || respTopStreamer.StatusCode() != 200 {
		return nil, fmt.Errorf("top Stream fetch Err: %v", respTopStreamer.Error())
	}
	return topStreamers.Data, nil
}

func (t *TwitchHandle) GetToken() (string, time.Time, error) {
	tokenResp := &TwitchToken{}
	resp, err := client.R().SetQueryParams(map[string]string{
		"client_id":     t.client_id,
		"client_secret": t.client_secret,
		"grant_type":    "client_credentials",
	}).SetResult(tokenResp).Post("https://id.twitch.tv/oauth2/token")

	if err != nil || resp.StatusCode() != 200 || tokenResp.AccessToken == "" {
		log.Printf("Err fetch token %v", err)
		return "", time.Time{}, fmt.Errorf("err fetch token")
	}

	safeExpirationTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-10 * time.Hour)

	return tokenResp.AccessToken, safeExpirationTime, nil
}

func (t *TwitchHandle) GetValidToken() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !tokenExpires.IsZero() && time.Now().Before(tokenExpires) {
		log.Printf("Token is still valid: %v", token)
		return token, nil
	}

	newToken, expirationTime, err := t.GetToken()
	if err != nil {
		return "", fmt.Errorf("cant fetch new token: %v", err)
	}

	token = newToken
	tokenExpires = expirationTime

	return token, nil
}
