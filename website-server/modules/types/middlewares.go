package types

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	OnlyAuthorized() fiber.Handler
	OnlyUnAuthorized() fiber.Handler
}
