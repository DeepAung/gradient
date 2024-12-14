package utils

import (
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetPayload(c *fiber.Ctx) (types.Payload, bool) {
	token, ok := c.Locals("token").(*jwt.Token)
	if !ok {
		return types.Payload{}, false
	}

	claims, ok := token.Claims.(*types.JwtClaims)
	if !ok {
		return types.Payload{}, false
	}

	return claims.Payload, true
}
