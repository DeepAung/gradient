package auth

import (
	"errors"
	"strconv"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/DeepAung/gradient/public-server/views/components"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrTokenIdNotFound        = errors.New("token id not found")
	ErrInvalidConfirmPassword = errors.New("password is not the same")
)

type AuthHandler struct {
	svc types.AuthSvc
	cfg *config.Config
}

func InitAuthHandler(router fiber.Router, svc types.AuthSvc, cfg *config.Config) {
	handler := &AuthHandler{
		svc: svc,
		cfg: cfg,
	}

	router.Post("/signin", handler.SignIn)
	router.Post("/signup", handler.SignUp)
	router.Post("/signout", handler.SignOut)
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var dto types.SignInDTO
	if err := c.BodyParser(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	if err := utils.Validate(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}

	passport, err := h.svc.SignIn(dto.Email, dto.Password)
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}

	utils.SetCookie(c, "accessToken", passport.Token.AccessToken, h.cfg.Jwt.AccessExpires)
	utils.SetCookie(c, "refreshToken", passport.Token.RefreshToken, h.cfg.Jwt.RefreshExpires)
	utils.SetCookie(c, "tokenId", strconv.Itoa(passport.Token.Id), h.cfg.Jwt.RefreshExpires)

	c.Response().Header.Add("HX-Redirect", "/home")
	return nil
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var dto types.SignUpDTO
	if err := c.BodyParser(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	if err := utils.Validate(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	if dto.Password != dto.ConfirmPassword {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(ErrInvalidConfirmPassword.Error())
	}

	passport, err := h.svc.SignUp(dto.Username, dto.Email, dto.Password)
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}

	utils.SetCookie(c, "accessToken", passport.Token.AccessToken, h.cfg.Jwt.AccessExpires)
	utils.SetCookie(c, "refreshToken", passport.Token.RefreshToken, h.cfg.Jwt.RefreshExpires)
	utils.SetCookie(c, "tokenId", strconv.Itoa(passport.Token.Id), h.cfg.Jwt.RefreshExpires)

	c.Response().Header.Add("HX-Redirect", "/home")
	return nil
}

func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	tokenIdStr := c.Cookies("tokenId", "")
	if tokenIdStr == "" {
		return c.SendString(ErrTokenIdNotFound.Error())
	}
	tokenId, err := strconv.Atoi(tokenIdStr)
	if err != nil {
		return c.SendString(err.Error())
	}

	if err := h.svc.SignOut(tokenId); err != nil {
		c.Response().Header.Add("HX-Retarget", "#alert")
		return utils.Render(c, components.AlertError(err.Error()))
	}

	utils.DeleteCookie(c, "accessToken")
	utils.DeleteCookie(c, "refreshToken")
	utils.DeleteCookie(c, "tokenId")

	c.Response().Header.Add("HX-Redirect", "/signin")
	return nil
}
