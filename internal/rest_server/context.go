package rest_server

import "context"

type contextKey int

const (
	contextUserID contextKey = iota
)

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextUserID, userID)
}

func getUserIdFromContext(ctx context.Context) string {
	return ctx.Value(contextUserID).(string)
}
