package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// NewStorage's redis branch requires a live Redis server, and is
// already exercised indirectly by the shared storage_test.go suite
// when Redis is available. Here we only cover the default,
// no-dependency branch.
func Test__NewStorage__DefaultsToInMemory(t *testing.T) {
	t.Run("DECK_STORAGE_TYPE unset", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "")

		s, err := NewStorage()
		require.NoError(t, err)
		require.NotNil(t, s)
		require.IsType(t, &InMemoryStorage{}, s)
	})

	t.Run("DECK_STORAGE_TYPE with unknown value", func(t *testing.T) {
		t.Setenv("DECK_STORAGE_TYPE", "something-unknown")

		s, err := NewStorage()
		require.NoError(t, err)
		require.NotNil(t, s)
		require.IsType(t, &InMemoryStorage{}, s)
	})
}
