package users

import (
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/pkg/utils"
	"github.com/DeepAung/gradient/website-server/views/components"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type usersHandler struct {
	usersSvc types.UsersSvc
}

func InitUsersHandler(
	router fiber.Router,
	mid types.Middleware,
	usersSvc types.UsersSvc,
) {
	handler := &usersHandler{
		usersSvc: usersSvc,
	}

	usersGroup := router.Group("/users")
	usersGroup.Put("/", mid.OnlyAuthorized(), handler.UpdateUser)
	usersGroup.Delete("/", mid.OnlyAuthorized(), handler.DeleteUser)
}

func (h *usersHandler) UpdateUser(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	var dto types.UpdateUserReq
	if err := c.BodyParser(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	if err := utils.Validate(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	if dto.CurrentPassword != nil && *dto.CurrentPassword == "" {
		dto.CurrentPassword = nil
	}
	if dto.NewPassword != nil && *dto.NewPassword == "" {
		dto.NewPassword = nil
	}

	if (dto.CurrentPassword == nil) && (dto.NewPassword != nil) {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString("empty current password")
	}
	if (dto.CurrentPassword != nil) && (dto.NewPassword == nil) {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString("empty new password")
	}

	picture, err := c.FormFile("picture")
	if err != nil {
		if err.Error() != fasthttp.ErrMissingFile.Error() {
			c.Response().Header.Add("HX-Retarget", "#error-text")
			return c.SendString(err.Error())
		}
	} else {
		newUrl, err := h.usersSvc.ReplacePicture(payload.UserId, payload.Email, picture)
		if err != nil {
			c.Response().Header.Add("HX-Retarget", "#error-text")
			_, msg := utils.ParseError(err)
			return c.SendString(msg)
		}

		dto.PictureUrl = &newUrl
	}

	if err := h.usersSvc.UpdateUser(payload.UserId, dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		_, msg := utils.ParseError(err)
		return c.SendString(msg)
	}

	return utils.RenderAlert(c, templ.Join(
		components.AlertSuccess("user information updated"),
		components.OOBWrap("innerHTML:#error-text", components.Text("")),
	))
}

func (h *usersHandler) DeleteUser(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		return c.Redirect("/signin", fiber.StatusFound)
	}

	if err := h.usersSvc.DeleteUser(payload.UserId); err != nil {
		_, msg := utils.ParseError(err)
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(msg)
	}

	c.Response().Header.Add("HX-Redirect", "/signin")
	return nil
}
