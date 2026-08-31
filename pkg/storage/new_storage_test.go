package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__NewStorage(t *testing.T) {
	t.Run("no DECK_STORAGE_TYPE set -> in-memory storage", func(t *testing.T) {
		s, err := NewStorage()
		require.NoError(t, err)
		require.IsType(t, &InMemoryStorage{}, s)
	})

	t.Run("DECK_STORAGE_TYPE=redis -> redis storage", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "redis")
		t.Setenv("REDIS_HOST", "localhost")
		t.Setenv("REDIS_PORT", "6379")

		s, err := NewStorage()
		require.NoError(t, err)
		require.IsType(t, &RedisStorage{}, s)
	})

	t.Run("DECK_STORAGE_TYPE=redis without REDIS_HOST -> error", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "redis")
		t.Setenv("REDIS_HOST", "")

		_, err := NewStorage()
		require.Error(t, err)
	})
}
