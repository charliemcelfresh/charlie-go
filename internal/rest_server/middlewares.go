package rest_server

import (
	"net/http"
)

const (
	XUserID         = "X-User-Id"
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

func AddContentTypeToResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(contentType, applicationJSON)
			next.ServeHTTP(w, r)
		},
	)
}

func AddUserIDToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get(XUserID)
			if userID != "" {
				ctx := r.Context()
				ctx = setUserIDInContext(ctx, userID)
				r = r.WithContext(ctx)
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Forbidden"))
			}
			next.ServeHTTP(w, r)
		},
	)
}
