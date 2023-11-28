package twirp_server

import (
	"context"
	"net/http"
)

const (
	Authorization = "Authorization"
)

func AddJwtTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get(Authorization)
			if auth != "" {
				ctx := r.Context()
				ctx = context.WithValue(ctx, contextJWT, auth)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		},
	)
}
