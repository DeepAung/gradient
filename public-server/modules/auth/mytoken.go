package auth

import (
	"errors"
	"time"

	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/golang-jwt/jwt/v5"
)

type tokenType string

const (
	accessTokenType  tokenType = "access-token"
	refreshTokenType tokenType = "refresh-token"
)

func generateToken(
	tokenType tokenType,
	duration time.Duration,
	secretKey []byte,
	payload types.Payload,
) (string, error) {
	claims := types.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gradient-server",
			Subject:   string(tokenType),
			Audience:  []string{"users", "admin"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Payload: payload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func parseToken(tokenString string, secretKey []byte) (*types.JwtClaims, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &types.JwtClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.JwtClaims); !ok {
		return nil, errors.New("invalid claims type")
	} else {
		return claims, nil
	}
}
