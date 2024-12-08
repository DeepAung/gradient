package authsvc

import (
	"errors"
	"fmt"

	"github.com/DeepAung/gradient/public-server/api/auth"
	"github.com/DeepAung/gradient/public-server/api/auth/authrepo"
	"github.com/DeepAung/gradient/public-server/api/users"
	"github.com/DeepAung/gradient/public-server/api/users/usersrepo"
	"github.com/DeepAung/gradient/public-server/pkg/config"
	"github.com/DeepAung/gradient/public-server/pkg/mytoken"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
)

type AuthSvc struct {
	authRepo  *authrepo.AuthRepo
	usersRepo *usersrepo.UsersRepo
	cfg       *config.Config
}

func NewAuthSvc(
	authRepo *authrepo.AuthRepo,
	usersRepo *usersrepo.UsersRepo,
	cfg *config.Config,
) *AuthSvc {
	return &AuthSvc{
		authRepo:  authRepo,
		usersRepo: usersRepo,
		cfg:       cfg,
	}
}

func (s *AuthSvc) SignUp(username, email, password string) (auth.Passport, error) {
	hashedPassword, err := utils.Hash(password)
	if err != nil {
		return auth.Passport{}, err
	}

	user, err := s.usersRepo.CreateUser(username, email, hashedPassword)
	return s.generatePassport(user)
}

func (s *AuthSvc) SignIn(email, password string) (auth.Passport, error) {
	if password == "" {
		return auth.Passport{}, ErrInvalidEmailOrPassword
	}

	user, err := s.usersRepo.FindOneUserWithPasswordByEmail(email)
	if err != nil {
		return auth.Passport{}, err
	}

	if !utils.Compare(password, user.Password) {
		return auth.Passport{}, ErrInvalidEmailOrPassword
	}

	return s.generatePassport(users.User{
		Id:         user.Id,
		Username:   user.Username,
		Email:      email,
		PictureUrl: user.PictureUrl,
		IsAdmin:    user.IsAdmin,
	})
}

func (s *AuthSvc) SignOut(tokenId int) error {
	return s.authRepo.DeleteToken(tokenId)
}

func (s *AuthSvc) UpdateTokens(tokenId int, refreshToken string) (auth.Token, error) {
	has, err := s.authRepo.HasToken(tokenId, refreshToken)
	if err != nil {
		return auth.Token{}, err
	}
	if !has {
		return auth.Token{}, ErrInvalidRefreshToken
	}

	claims, err := mytoken.ParseToken(mytoken.RefreshToken, refreshToken, s.cfg.Jwt.SecretKey)
	if err != nil {
		return auth.Token{}, err
	}

	newAccessToken, err := mytoken.GenerateToken(
		mytoken.AccessToken,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return auth.Token{}, err
	}

	newRefreshToken, err := mytoken.GenerateToken(
		mytoken.RefreshToken,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return auth.Token{}, err
	}

	if err := s.authRepo.UpdateTokens(tokenId, newAccessToken, newRefreshToken); err != nil {
		return auth.Token{}, err
	}

	return auth.Token{
		Id:           tokenId,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func keyFunc(key []byte) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return key, nil
	}
}

func (s *AuthSvc) generatePassport(user users.User) (auth.Passport, error) {
	payload := auth.Payload{
		UserId:   user.Id,
		Username: user.Username,
	}
	accessToken, err := mytoken.GenerateToken(
		mytoken.AccessToken,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return auth.Passport{}, err
	}
	refreshToken, err := mytoken.GenerateToken(
		mytoken.RefreshToken,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return auth.Passport{}, err
	}

	token, err := s.authRepo.CreateToken(accessToken, refreshToken)
	if err != nil {
		return auth.Passport{}, nil
	}

	return auth.Passport{
		User:  user,
		Token: token,
	}, nil
}
