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
	client       *resty.Client
	clientID     string
	clientSecret string
	mu           sync.Mutex
	token        string
	tokenExpires time.Time
}

func NewTwitchHandle(cfg *config.EnvVars) *TwitchHandle {
	return &TwitchHandle{
		client:       resty.New(),
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
	}
}

func (t *TwitchHandle) GetTopGames() ([]Game, error) {
	var authHeaders = map[string]string{
		"Client-ID":     t.clientID,
		"Authorization": fmt.Sprintf("Bearer %s", t.token),
	}
	topGames := &TopGamesResponse{}
	respTopGames, err := t.client.R().SetHeaders(authHeaders).SetQueryParam("first", strconv.Itoa(50)).SetResult(topGames).Get("https://api.twitch.tv/helix/games/top")

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

var ErrFetchTopStream = errors.New("error fetch Top Stream")

func (t *TwitchHandle) GetTopStream(gameId string) ([]Stream, error) {
	var authHeaders = map[string]string{
		"Client-ID":     t.clientID,
		"Authorization": fmt.Sprintf("Bearer %s", t.token),
	}

	var topStreamers = &StreamsResponse{}
	respTopStreamer, err := t.client.R().SetHeaders(authHeaders).SetQueryParams(map[string]string{
		"game_id": gameId,
		"first":   strconv.Itoa(100),
	}).SetResult(topStreamers).Get("https://api.twitch.tv/helix/streams")

	if err != nil || respTopStreamer.StatusCode() != 200 {
		log.Printf("err fetch top stream %v", respTopStreamer.Error())
		return nil, ErrFetchTopStream
	}
	return topStreamers.Data, nil
}

var ErrFetchToken = errors.New("Error fetch token")

func (t *TwitchHandle) GetToken() (string, time.Time, error) {
	tokenResp := &TwitchToken{}
	resp, err := t.client.R().SetQueryParams(map[string]string{
		"client_id":     t.clientID,
		"client_secret": t.clientSecret,
		"grant_type":    "client_credentials",
	}).SetResult(tokenResp).Post("https://id.twitch.tv/oauth2/token")

	if err != nil || resp.StatusCode() != 200 || tokenResp.AccessToken == "" {
		log.Printf("Err fetch token %v", err)
		return "", time.Time{}, ErrFetchToken
	}
	log.Printf("GET TIME TOKEN %v", tokenResp.ExpiresIn)
	safeExpirationTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-1 * time.Minute)

	t.token = tokenResp.AccessToken
	t.tokenExpires = safeExpirationTime

	return tokenResp.AccessToken, safeExpirationTime, nil
}

// var ErrTokenIsValid = errors.New("error fetch token")
var ErrFetchNewToken = errors.New("error fetch new token")

func (t *TwitchHandle) GetValidToken() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.tokenExpires.IsZero() && time.Now().Before(t.tokenExpires) {
		log.Printf("Token is still valid: %v", t.token)
		return t.token, nil
	}

	newToken, expirationTime, err := t.GetToken()
	if err != nil {
		return "", ErrFetchNewToken
	}

	t.token = newToken
	t.tokenExpires = expirationTime

	return t.token, nil
}
