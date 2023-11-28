package jwt_maker

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func Run(jwtDuration string, userID string) {
	dur, err := time.ParseDuration(jwtDuration)
	if err != nil {
		panic(err)
	}
	prvKey, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}
	pubKey, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		panic(err)
	}

	jwtToken := NewJWT(prvKey, pubKey)

	token, err := jwtToken.Create(dur, "charlie-go", userID)
	if err != nil {
		panic(err)
	}

	_, err = jwtToken.Validate(token)
	if err != nil {
		panic(err)
	}
	// Print token to STDOUT
	fmt.Println(token)
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) Create(ttl time.Duration, iss string, userID string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["iss"] = iss                 // issuing service
	claims["aud"] = "user"              // audience, eg user, admin, service-to-service
	claims["sub"] = userID
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	// Print token to STDOUT so we can use it in requests
	return token, nil
}

func (j JWT) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
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
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
