package utils

import (
	"strconv"
	"time"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/gofiber/fiber/v2"
)

func SetTokenCookies(c *fiber.Ctx, token types.Token, cfg *config.Config) {
	SetCookie(c, "accessToken", token.AccessToken, cfg.Jwt.AccessExpires)
	SetCookie(c, "refreshToken", token.RefreshToken, cfg.Jwt.RefreshExpires)
	SetCookie(c, "tokenId", strconv.Itoa(token.Id), cfg.Jwt.RefreshExpires)
}

func DeleteTokenCookies(c *fiber.Ctx) {
	DeleteCookie(c, "accessToken")
	DeleteCookie(c, "refreshToken")
	DeleteCookie(c, "tokenId")
}

func SetCookie(c *fiber.Ctx, name, value string, maxAge time.Duration) {
	var expires time.Time
	if maxAge != 0 {
		expires = time.Now().Add(maxAge)
	}

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(maxAge.Seconds()),
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func DeleteCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-7 * 24 * time.Hour),
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}
