package middlewares

import (
	"strconv"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/pkg/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type middlewareImpl struct {
	cfg     *config.Config
	authSvc types.AuthSvc
}

func NewMiddleware(cfg *config.Config, authSvc types.AuthSvc) types.Middleware {
	return &middlewareImpl{
		cfg:     cfg,
		authSvc: authSvc,
	}
}

func (m *middlewareImpl) OnlyAuthorized() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: m.cfg.Jwt.SecretKey},
		TokenLookup: "cookie:accessToken",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Try update token
			refreshToken := c.Cookies("refreshToken")
			tokenId, err := strconv.Atoi(c.Cookies("tokenId"))
			if err != nil {
				utils.DeleteTokenCookies(c)
				return c.Redirect("/signin", fiber.StatusFound)
			}

			token, claims, err := m.authSvc.UpdateTokens(tokenId, refreshToken)
			if err != nil {
				return c.Redirect("/signin", fiber.StatusFound)
			}

			utils.SetTokenCookies(c, token, m.cfg)
			c.Locals("token", &jwt.Token{Claims: claims})

			return c.Next()
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ContextKey: "token",
		Claims:     &types.JwtClaims{},
	})
}

func (m *middlewareImpl) OnlyUnAuthorized() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: m.cfg.Jwt.SecretKey},
		TokenLookup: "cookie:accessToken",
		SuccessHandler: func(c *fiber.Ctx) error {
			redirectTo := c.Query("redirect", "/home")
			return c.Redirect(redirectTo, fiber.StatusFound)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next()
		},
	})
}
