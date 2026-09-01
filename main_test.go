package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__GetPort(t *testing.T) {
	t.Run("unset API_PORT -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "")
		require.Equal(t, defaultPort, getPort())
	})

	t.Run("valid API_PORT -> that port", func(t *testing.T) {
		t.Setenv("API_PORT", "8080")
		require.Equal(t, 8080, getPort())
	})

	t.Run("invalid API_PORT -> default port", func(t *testing.T) {
		t.Setenv("API_PORT", "not-a-number")
		require.Equal(t, defaultPort, getPort())
	})
}
