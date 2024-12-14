package auth

import (
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidEmailOrPassword = fiber.NewError(
		fiber.StatusBadRequest,
		"invalid email or password",
	)
	ErrInvalidRefreshTokenOrTokenNotFound = fiber.NewError(fiber.StatusBadRequest,
		"invalid refresh token or token not found",
	)
)

type authSvc struct {
	authRepo  types.AuthRepo
	usersRepo types.UsersRepo
	cfg       *config.Config
}

func NewAuthSvc(
	authRepo types.AuthRepo,
	usersRepo types.UsersRepo,
	cfg *config.Config,
) types.AuthSvc {
	return &authSvc{
		authRepo:  authRepo,
		usersRepo: usersRepo,
		cfg:       cfg,
	}
}

func (s *authSvc) SignUp(username, email, password string) (types.Passport, error) {
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

func (s *authSvc) SignIn(email, password string) (types.Passport, error) {
	if password == "" {
		return types.Passport{}, ErrInvalidEmailOrPassword
	}

	user, err := s.usersRepo.FindOneUserWithPasswordByEmail(email)
	if err != nil {
		return types.Passport{}, err
	}

	if !utils.CompareHash(password, user.Password) {
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

func (s *authSvc) SignOut(tokenId int) error {
	return s.authRepo.DeleteToken(tokenId)
}

func (s *authSvc) UpdateTokens(
	tokenId int,
	refreshToken string,
) (types.Token, *types.JwtClaims, error) {
	has, err := s.authRepo.HasToken(tokenId, refreshToken)
	if err != nil {
		return types.Token{}, nil, err
	}
	if !has {
		return types.Token{}, nil, ErrInvalidRefreshTokenOrTokenNotFound
	}

	claims, err := parseToken(refreshToken, s.cfg.Jwt.SecretKey)
	if err != nil {
		return types.Token{}, nil, err
	}

	newAccessToken, err := generateToken(
		accessTokenType,
		s.cfg.Jwt.AccessExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, nil, err
	}

	newRefreshToken, err := generateToken(
		refreshTokenType,
		s.cfg.Jwt.RefreshExpires,
		s.cfg.Jwt.SecretKey,
		claims.Payload,
	)
	if err != nil {
		return types.Token{}, nil, err
	}

	if err := s.authRepo.UpdateTokens(tokenId, newAccessToken, newRefreshToken); err != nil {
		return types.Token{}, nil, err
	}

	return types.Token{
		Id:           tokenId,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, claims, nil
}

func (s *authSvc) generatePassport(user types.User) (types.Passport, error) {
	payload := types.Payload{
		UserId: user.Id,
		Email:  user.Email,
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
