package auth

import (
	"github.com/DeepAung/gradient/public-server/api/users"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Id           int    `db:"id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
}

type Passport struct {
	User  users.User
	Token Token
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Payload Payload
}

type Payload struct {
	UserId   int
	Username string
}
