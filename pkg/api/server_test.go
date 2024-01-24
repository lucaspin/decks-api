package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__HealthCheckEndpointRespondsWith200(t *testing.T) {
	testServer := NewServer()
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	testServer.router.ServeHTTP(response, request)
	require.Equal(t, response.Code, 200)
}
