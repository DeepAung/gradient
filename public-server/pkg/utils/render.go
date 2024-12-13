package utils

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RenderAlert(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.Add("HX-Retarget", "#alerts")
	c.Response().Header.Add("HX-Reswap", "beforeend")
	return Render(c, component)
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Context().SetContentType("text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func ParseError(err error) (status int, msg string) {
	switch err := err.(type) {
	case *fiber.Error:
		status = err.Code
		msg = err.Message
	default:
		tmp := fiber.NewError(fiber.StatusInternalServerError)
		status = tmp.Code
		msg = tmp.Message
	}

	if status >= 500 {
		log.Error(err.Error())
	}

	return
}
