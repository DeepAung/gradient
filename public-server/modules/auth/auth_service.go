package auth

import (
	"errors"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
)

var (
	ErrInvalidEmailOrPassword             = errors.New("invalid email or password")
	ErrInvalidRefreshTokenOrTokenNotFound = errors.New("invalid refresh token or token not found")
)

type AuthSvc struct {
	authRepo  types.AuthRepo
	usersRepo types.UsersRepo
	cfg       *config.Config
}

func NewAuthSvc(
	authRepo types.AuthRepo,
	usersRepo types.UsersRepo,
	cfg *config.Config,
) types.AuthSvc {
	return &AuthSvc{
		authRepo:  authRepo,
		usersRepo: usersRepo,
		cfg:       cfg,
	}
}

func (s *AuthSvc) SignUp(username, email, password string) (types.Passport, error) {
	hashedPassword, err := utils.Hash(password)
	if err != nil {
		return types.Passport{}, err
	}

	user, err := s.usersRepo.CreateUser(username, email, hashedPassword)
	if err != nil {
		return types.Passport{}, err
	}
	return s.generatePassport(user)
}

func (s *AuthSvc) SignIn(email, password string) (types.Passport, error) {
	if password == "" {
		return types.Passport{}, ErrInvalidEmailOrPassword
	}

	user, err := s.usersRepo.FindOneUserWithPasswordByEmail(email)
	if err != nil {
		return types.Passport{}, err
	}

	if !utils.Compare(password, user.Password) {
		return types.Passport{}, ErrInvalidEmailOrPassword
	}

	return s.generatePassport(types.User{
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

func (s *AuthSvc) UpdateTokens(tokenId int, refreshToken string) (types.Token, error) {
	has, err := s.authRepo.HasToken(tokenId, refreshToken)
	if err != nil {
		return types.Token{}, err
	}
	if !has {
		return types.Token{}, ErrInvalidRefreshTokenOrTokenNotFound
	}

	claims, err := parseToken(refreshToken, s.cfg.Jwt.SecretKey)
	if err != nil {
		return types.Token{}, err
	}

	newAccessToken, err := generateToken(
		accessTokenType,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, err
	}

	newRefreshToken, err := generateToken(
		refreshTokenType,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, err
	}

	if err := s.authRepo.UpdateTokens(tokenId, newAccessToken, newRefreshToken); err != nil {
		return types.Token{}, err
	}

	return types.Token{
		Id:           tokenId,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthSvc) generatePassport(user types.User) (types.Passport, error) {
	payload := types.Payload{
		UserId:   user.Id,
		Username: user.Username,
	}
	accessToken, err := generateToken(
		accessTokenType,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return types.Passport{}, err
	}
	refreshToken, err := generateToken(
		refreshTokenType,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		payload,
	)
	if err != nil {
		return types.Passport{}, err
	}

	token, err := s.authRepo.CreateToken(accessToken, refreshToken)
	if err != nil {
		return types.Passport{}, nil
	}

	return types.Passport{
		User:  user,
		Token: token,
	}, nil
}
