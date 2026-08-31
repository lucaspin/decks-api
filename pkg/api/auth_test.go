package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__AuthMiddleware_PassesRequestsThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("ok"))
	})

	handler := authMiddleware(next)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.True(t, called, "expected the wrapped handler to be called")
	require.Equal(t, http.StatusTeapot, rr.Code)
	require.Equal(t, "ok", rr.Body.String())
}
