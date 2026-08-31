package storage

import (
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
		require.EqualError(t, err, "no REDIS_HOST set")
	})

	t.Run("REDIS_PORT unset -> defaults to 6379", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "my-redis-host")
		t.Setenv("REDIS_PORT", "")
		t.Setenv("REDIS_USERNAME", "")
		t.Setenv("REDIS_PASSWORD", "")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "my-redis-host",
			Port:     "6379",
			Username: "",
			Password: "",
		}, config)
	})

	t.Run("all fields set from environment", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "my-redis-host")
		t.Setenv("REDIS_PORT", "6380")
		t.Setenv("REDIS_USERNAME", "my-user")
		t.Setenv("REDIS_PASSWORD", "my-password")

		config, err := NewRedisConfigFromEnvironment()
		require.NoError(t, err)
		require.Equal(t, &RedisConfig{
			Host:     "my-redis-host",
			Port:     "6380",
			Username: "my-user",
			Password: "my-password",
		}, config)
	})
}
