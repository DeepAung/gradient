package types

import "github.com/golang-jwt/jwt/v5"

type AuthSvc interface {
	SignUp(username, email, password string) (Passport, error)
	SignIn(email, password string) (Passport, error)
	SignOut(tokenId int) error
	UpdateTokens(tokenId int, refreshToken string) (Token, error)
}

type AuthRepo interface {
	HasToken(id int, refreshToken string) (bool, error)
	CreateToken(accessToken, refreshToken string) (Token, error)
	UpdateTokens(id int, newAccessToken, newRefreshToken string) error
	DeleteToken(tokenId int) error
}

type Token struct {
	Id           int    `db:"id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
}

type Passport struct {
	User  User
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