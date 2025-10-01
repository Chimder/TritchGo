package handlers

import (
	"log"
	"testing"
	"time"
	"tritchgo/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestGetTopStream(t *testing.T) {}
func TestGetToken(t *testing.T) {
	cfg := config.LoadEnv()
	log.Printf("ID %s", cfg.ClientID)
	log.Printf("SECRET %s", cfg.ClientSecret)
	handlerTwitch := NewTwitchHandle(cfg)
	assert.NotNil(t, handlerTwitch)

	var token string
	var time time.Time
	var tErr error
	t.Run("Get token", func(t *testing.T) {
		token, time, tErr = handlerTwitch.GetToken()
		require.NoError(t, tErr)
		require.NotEmpty(t, token)
		require.NotEmpty(t, time)
	})

	log.Printf("TIME1 %v", time)
	t.Run("Get valid token", func(t *testing.T) {
		valid, err := handlerTwitch.GetValidToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, valid)
		assert.Equal(t, token, valid)
	})

	t.Run("Get top games", func(t *testing.T) {
		topGames, err := handlerTwitch.GetTopGames()
		assert.NoError(t, err)
		assert.NotNil(t, topGames)
		assert.Greater(t, len(topGames), 0)
		assert.IsType(t, []Game{}, topGames)

		first := topGames[0]
		assert.NotEmpty(t, first.ID)
		assert.NotEmpty(t, first.Name)
		t.Run("Get top stream", func(t *testing.T) {
			topStream, err := handlerTwitch.GetTopStream(first.ID)
			assert.NoError(t, err)
			assert.NotNil(t, topStream)
			assert.Greater(t, len(topStream), 0)
			assert.IsType(t, []Stream{}, topStream)

			firstStream := topStream[0]
			assert.NotEmpty(t, firstStream.GameName)
			assert.NotEmpty(t, firstStream.GameID)
		})
	})
}
