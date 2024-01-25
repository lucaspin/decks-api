package api

import (
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// NOTE: this is where authentication would be implemented.
		// There is no requirement about authentication on the task,
		// so I'm not going to do any kind of authentication at all.

		next.ServeHTTP(w, r)
	})
}
