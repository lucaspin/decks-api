package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__GetPort(t *testing.T) {
	t.Run("API_PORT not set -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "")
		require.Equal(t, defaultPort, getPort())
	})

	t.Run("API_PORT set to a valid number -> uses it", func(t *testing.T) {
		t.Setenv("API_PORT", "8012")
		require.Equal(t, 8012, getPort())
	})

	t.Run("API_PORT set to an invalid number -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "notanumber")
		require.Equal(t, defaultPort, getPort())
	})
}
