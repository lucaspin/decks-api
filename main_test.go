package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__getPort(t *testing.T) {
	t.Run("no API_PORT set -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "")
		require.Equal(t, defaultPort, getPort())
	})

	t.Run("invalid API_PORT set -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "not-a-number")
		require.Equal(t, defaultPort, getPort())
	})

	t.Run("valid API_PORT set -> parsed port", func(t *testing.T) {
		t.Setenv("API_PORT", "8012")
		require.Equal(t, 8012, getPort())
	})
}
