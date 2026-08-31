package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__NewRedisConfigFromEnvironment(t *testing.T) {
	t.Run("no REDIS_HOST set -> error", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "")

		config, err := NewRedisConfigFromEnvironment()
		require.Nil(t, config)
		require.ErrorContains(t, err, "no REDIS_HOST set")
	})

	t.Run("REDIS_HOST set, REDIS_PORT unset -> defaults port to 6379", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "redis.example.com")
		t.Setenv("REDIS_PORT", "")
		t.Setenv("REDIS_USERNAME", "user")
		t.Setenv("REDIS_PASSWORD", "pass")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "redis.example.com",
			Port:     "6379",
			Username: "user",
			Password: "pass",
		}, config)
	})

	t.Run("REDIS_HOST and REDIS_PORT set -> uses both", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "redis.example.com")
		t.Setenv("REDIS_PORT", "1234")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, "redis.example.com", config.Host)
		require.Equal(t, "1234", config.Port)
	})
}
