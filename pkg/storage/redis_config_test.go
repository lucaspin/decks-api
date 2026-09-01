package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__NewRedisConfigFromEnvironment(t *testing.T) {
	t.Run("missing REDIS_HOST -> error", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "")
		t.Setenv("REDIS_PORT", "")
		t.Setenv("REDIS_USERNAME", "")
		t.Setenv("REDIS_PASSWORD", "")

		config, err := NewRedisConfigFromEnvironment()
		require.Nil(t, config)
		require.ErrorContains(t, err, "no REDIS_HOST set")
	})

	t.Run("host set with no port -> defaults to 6379", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "myredis")
		t.Setenv("REDIS_PORT", "")
		t.Setenv("REDIS_USERNAME", "")
		t.Setenv("REDIS_PASSWORD", "")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "myredis",
			Port:     "6379",
			Username: "",
			Password: "",
		}, config)
	})

	t.Run("explicit port/username/password are carried through", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "myredis")
		t.Setenv("REDIS_PORT", "1234")
		t.Setenv("REDIS_USERNAME", "user")
		t.Setenv("REDIS_PASSWORD", "pass")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "myredis",
			Port:     "1234",
			Username: "user",
			Password: "pass",
		}, config)
	})
}

func Test__NewStorage(t *testing.T) {
	t.Run("unset DECK_STORAGE_TYPE -> in-memory storage", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "")

		store, err := NewStorage()
		require.NoError(t, err)
		require.IsType(t, &InMemoryStorage{}, store)
	})

	t.Run("unrecognized DECK_STORAGE_TYPE -> in-memory storage", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "something-else")

		store, err := NewStorage()
		require.NoError(t, err)
		require.IsType(t, &InMemoryStorage{}, store)
	})

	t.Run("redis DECK_STORAGE_TYPE with valid environment -> redis storage", func(t *testing.T) {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			t.Skip("REDIS_HOST not set - skipping redis-backed NewStorage test")
		}

		t.Setenv("DECK_STORAGE_TYPE", "redis")

		store, err := NewStorage()
		require.NoError(t, err)
		require.IsType(t, &RedisStorage{}, store)
	})
}
