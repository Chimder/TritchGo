package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	redistest "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redistest.Run(ctx, "redis:latest")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = redisContainer.Terminate(ctx)
	})

	connStr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)
	connStr = strings.TrimPrefix(connStr, "redis://")

	client := RedisDb(connStr)
	require.NotNil(t, client, "redis is nil")

	status, err := client.Ping(ctx).Result()
	require.NoError(t, err)
	require.Equal(t, "PONG", status)

	key := "test-key"
	val := "hello world"

	err = client.Set(ctx, key, val, 10*time.Second).Err()
	require.NoError(t, err)

	got, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	require.Equal(t, val, got)
}
