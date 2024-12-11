package middlewares

import (
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/types"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	cfg *config.Config
}

func NewMiddleware(cfg *config.Config) types.Middleware {
	return &Middleware{
		cfg: cfg,
	}
}

func (m *Middleware) OnlyAuthorized() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: m.cfg.Jwt.SecretKey},
		TokenLookup: "cookie:accessToken",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO:
			// try update tokens if error then
			return c.Redirect("/signin", fiber.StatusFound)
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ContextKey: "claims",
		Claims:     &types.JwtClaims{},
	})
}

func (m *Middleware) OnlyUnAuthorized() fiber.Handler {
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
