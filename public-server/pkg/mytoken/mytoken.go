package mytoken

import (
	"errors"
	"time"

	"github.com/DeepAung/gradient/public-server/api/auth"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access-token"
	RefreshToken TokenType = "refresh-token"
)

func ParseToken(
	tokenType TokenType,
	tokenString string,
	secretKey []byte,
) (*auth.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.JwtClaims{}, keyFunc(secretKey))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*auth.JwtClaims); !ok {
		return nil, errors.New("invalid claims type")
	} else {
		return claims, nil
	}
}

func keyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return key, nil
	}
}

func GenerateToken(
	tokenType TokenType,
	duration time.Duration,
	secretKey []byte,
	payload auth.Payload,
) (string, error) {
	claims := auth.JwtClaims{
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
