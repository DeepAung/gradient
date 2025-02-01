package views

import (
	"strconv"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/pkg/utils"
	"github.com/DeepAung/gradient/website-server/views/pages"
	"github.com/gofiber/fiber/v2"
)

type viewsHandler struct {
	usersSvc  types.UsersSvc
	tasksSvc  types.TasksSvc
	graderCfg *graderconfig.Config
}

func InitViewsHandler(
	router fiber.Router,
	mid types.Middleware,
	usersSvc types.UsersSvc,
	tasksSvc types.TasksSvc,
	graderCfg *graderconfig.Config,
) {
	handler := &viewsHandler{
		usersSvc:  usersSvc,
		tasksSvc:  tasksSvc,
		graderCfg: graderCfg,
	}

	router.Get("/", handler.Welcome)
	router.Get("/signup", mid.OnlyUnAuthorized(), handler.SignUp)
	router.Get("/signin", mid.OnlyUnAuthorized(), handler.SignIn)
	router.Get("/home", mid.OnlyAuthorized(), handler.Home)
	router.Get("/profile", mid.OnlyAuthorized(), handler.Profile)
	router.Get("/tasks/:id", mid.OnlyAuthorized(), handler.TaskDetail)
}

func (h *viewsHandler) Welcome(c *fiber.Ctx) error {
	return utils.Render(c, pages.Welcome())
}

func (h *viewsHandler) SignIn(c *fiber.Ctx) error {
	return utils.Render(c, pages.SignIn())
}

func (h *viewsHandler) SignUp(c *fiber.Ctx) error {
	return utils.Render(c, pages.SignUp())
}

func (h *viewsHandler) Home(c *fiber.Ctx) error {
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

func (h *viewsHandler) Profile(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	user, err := h.usersSvc.GetUser(payload.UserId)
	if err != nil {
		_, msg := utils.ParseError(err)
		return utils.Render(c, pages.Error(msg, "/home"))
	}

	return utils.Render(c, pages.Profile(user))
}

func (h *viewsHandler) TaskDetail(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	user, err := h.usersSvc.GetUser(payload.UserId)
	if err != nil {
		_, msg := utils.ParseError(err)
		return utils.Render(c, pages.Error(msg, "/home"))
	}

	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.Render(c, pages.Error("Invalid task id", "/home"))
	}

	task, err := h.tasksSvc.GetTask(payload.UserId, taskId)
	if err != nil {
		_, msg := utils.ParseError(err)
		return utils.Render(c, pages.Error(msg, "/home"))
	}

	languages := h.graderCfg.Languages
	return utils.Render(c, pages.TaskDetail(user, task, languages))
}
