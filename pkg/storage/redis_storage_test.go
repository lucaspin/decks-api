package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests only cover NewRedisConfigFromEnvironment, which is pure
// environment-variable parsing and does not require a live Redis server.
func Test__NewRedisConfigFromEnvironment(t *testing.T) {
	t.Run("missing REDIS_HOST -> error", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "")

		config, err := NewRedisConfigFromEnvironment()
		require.Nil(t, config)
		require.ErrorContains(t, err, "no REDIS_HOST set")
	})

	t.Run("REDIS_PORT unset -> defaults to 6379", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "localhost")
		t.Setenv("REDIS_PORT", "")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{Host: "localhost", Port: "6379"}, config)
	})

	t.Run("custom REDIS_PORT is respected", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "localhost")
		t.Setenv("REDIS_PORT", "1234")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{Host: "localhost", Port: "1234"}, config)
	})

	t.Run("REDIS_USERNAME and REDIS_PASSWORD are passed through", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "localhost")
		t.Setenv("REDIS_PORT", "1234")
		t.Setenv("REDIS_USERNAME", "user")
		t.Setenv("REDIS_PASSWORD", "pass")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "localhost",
			Port:     "1234",
			Username: "user",
			Password: "pass",
		}, config)
	})
}
