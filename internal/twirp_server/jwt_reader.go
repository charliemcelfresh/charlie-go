package twirp_server

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
)

const (
	invalidToken    = "invalid token"
	missingExpClaim = "missing exp claim"
	missingIatClaim = "missing iat claim"
	missingIssClaim = "missing iss claim"
	missingAudClaim = "missing aud claim"
	wrongAudience   = "aud is not user"
	missingSubClaim = "missing sub claim"
)

type Jwt struct {
	publicKey []byte
}

type Claims struct {
	Exp int64
	Iat int64
	Iss string
	Aud string
	Sub string
}

func NewJWT(publicKey []byte) Jwt {
	return Jwt{
		publicKey: publicKey,
	}
}

func (j Jwt) Validate(token string) (Claims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing public key: %w", err)
	}

	tok, err := jwt.Parse(
		token, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
			}

			return key, nil
		},
	)

	if err != nil {
		return Claims{}, fmt.Errorf("validate: %w", err)
	}

	c, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return Claims{}, fmt.Errorf(invalidToken)
	}

	exp, ok := c["exp"].(float64)
	if !ok {
		return Claims{}, fmt.Errorf(missingExpClaim)
	}

	iat, ok := c["iat"].(float64)
	if !ok {
		return Claims{}, fmt.Errorf(missingIatClaim)
	}

	iss, ok := c["iss"].(string)
	if !ok {
		return Claims{}, fmt.Errorf(missingIssClaim)
	}

	aud, ok := c["aud"].(string)
	if !ok {
		return Claims{}, fmt.Errorf(missingAudClaim)
	}

	if aud != "user" {
		return Claims{}, fmt.Errorf(wrongAudience)
	}

	sub, ok := c["sub"].(string)
	if !ok {
		return Claims{}, fmt.Errorf(missingSubClaim)
	}

	claimsToReturn := Claims{
		Exp: int64(exp),
		Iat: int64(iat),
		Iss: iss,
		Aud: aud,
		Sub: sub,
	}

	return claimsToReturn, nil
}
