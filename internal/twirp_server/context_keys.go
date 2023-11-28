package twirp_server

import "context"

type contextKey int

const (
	contextUserID contextKey = iota
	contextJWT
)

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextUserID, userID)
}

func getUserIdFromContext(ctx context.Context) string {
	return ctx.Value(contextUserID).(string)
}
