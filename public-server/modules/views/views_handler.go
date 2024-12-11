package views

import (
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/DeepAung/gradient/public-server/views/pages"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct {
	usersSvc types.UsersSvc
	tasksSvc types.TasksSvc
}

func InitViewsHandler(
	router fiber.Router,
	mid types.Middleware,
	usersSvc types.UsersSvc,
	tasksSvc types.TasksSvc,
) {
	handler := &ViewsHandler{
		usersSvc: usersSvc,
		tasksSvc: tasksSvc,
	}

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
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	user, err := h.usersSvc.GetUser(payload.UserId)
	if err != nil {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	return utils.Render(c, pages.Home(user))
}
