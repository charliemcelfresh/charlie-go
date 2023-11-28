package twirp_server

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/twitchtv/twirp"
)

const (
	missingJWT       = "Missing JWT"
	missingPublicKey = "Missing public key"
)

func (p provider) AuthHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			claims, err := jWTIsValid(ctx)
			if err != nil {
				return context.Background(), twirp.NewError(twirp.Unauthenticated, err.Error())
			}
			ctx = setUserIDInContext(ctx, claims.Sub)
			return ctx, nil
		},
	}
}

func jWTIsValid(ctx context.Context) (Claims, error) {
	token, ok := ctx.Value(contextJWT).(string)
	if !ok {
		return Claims{}, errors.New(missingJWT)
	}
	token = strings.Split(token, "Bearer ")[1]
	pubKey, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		return Claims{}, errors.New(missingPublicKey)
	}
	jwtToken := NewJWT(pubKey)
	claims, err := jwtToken.Validate(token)
	if err != nil {
		return Claims{}, err
	}
	return claims, nil
}
