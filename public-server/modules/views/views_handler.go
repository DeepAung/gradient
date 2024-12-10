package views

import (
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/DeepAung/gradient/public-server/views/pages"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct{}

func InitViewsHandler(router fiber.Router, mid types.Middleware) {
	handler := &ViewsHandler{}
	router.Get("/", handler.Welcome)
	router.Get("/signup", mid.OnlyUnAuthorized(), handler.SignUp)
	router.Get("/signin", mid.OnlyUnAuthorized(), handler.SignIn)
	router.Get("/home", mid.OnlyAuthorized(), handler.Home)
}

func (h *ViewsHandler) Welcome(c *fiber.Ctx) error {
	return utils.Render(c, pages.Welcome())
}

func (h *ViewsHandler) SignIn(c *fiber.Ctx) error {
	return utils.Render(c, pages.SignIn())
}

func (h *ViewsHandler) SignUp(c *fiber.Ctx) error {
	return utils.Render(c, pages.SignUp())
}

func (h *ViewsHandler) Home(c *fiber.Ctx) error {
	return utils.Render(c, pages.Home(types.User{
		Id:         1,
		Username:   "DeepAung",
		Email:      "i.deepaung@gmail.com",
		PictureUrl: "",
		IsAdmin:    false,
	}))
}
