package submissions

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/hub"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/DeepAung/gradient/public-server/views/components"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type submissionsHandler struct {
	submissionsSvc types.SubmissionsSvc
	tasksSvc       types.TasksSvc
}

func InitSubmissionsHandler(
	router fiber.Router,
	mid types.Middleware,
	submissionsSvc types.SubmissionsSvc,
	tasksSvc types.TasksSvc,
) {
	handler := &submissionsHandler{
		submissionsSvc: submissionsSvc,
		tasksSvc:       tasksSvc,
	}

	submissionsGroup := router.Group("/submissions")

	submissionsGroup.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	submissionsGroup.Post("/", mid.OnlyAuthorized(), handler.SubmitCode)
	submissionsGroup.Get("/ws/connect", mid.OnlyAuthorized(), websocket.New(handler.SendResult))
}

func (h *submissionsHandler) SubmitCode(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		c.Response().Header.Add("HX-Redirect", "/signin")
		return nil
	}

	// Get task id
	taskId, err := strconv.Atoi(c.FormValue("task_id"))
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}

	// Get language
	language, ok := types.StringToProtoLanguage(c.FormValue("language"))
	if !ok {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(ErrInvalidLanguage.Error())
	}

	// Get code string
	codeFileHeader, err := c.FormFile("code_file")
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	codeFile, err := codeFileHeader.Open()
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	codeBytes, err := io.ReadAll(codeFile)
	if err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}
	codeFile.Close()

	dto := types.CreateSubmissionReq{
		UserId:   payload.UserId,
		TaskId:   taskId,
		Code:     string(codeBytes),
		Language: language,
	}

	if err := utils.Validate(&dto); err != nil {
		c.Response().Header.Add("HX-Retarget", "#error-text")
		return c.SendString(err.Error())
	}

	resultCh, createCh, err := h.submissionsSvc.SubmitCode(dto)
	if err != nil {
		_, msg := utils.ParseError(err)
		return c.SendString(msg)
	}

	resultId := uuid.NewString()
	utils.SetCookie(c, "resultId", resultId, 0)
	hub.CreateResult(resultId, resultCh)

	_ = createCh
	// createRes := <-createCh
	//
	// if createRes.Err != nil {
	// 	utils.DeleteCookie(c, "resultId")
	// 	hub.DeleteResult(resultId)
	// 	return c.SendString(createRes.Err.Error())
	// }

	c.Response().Header.Add("HX-Retarget", "#submission-results")
	return utils.Render(c, components.SubmissionResults())
}

func (h *submissionsHandler) SendResult(c *websocket.Conn) {
	resultId := c.Cookies("resultId")
	resultCh, ok := hub.PopResult(resultId)
	if !ok {
		c.Close()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for result := range resultCh {
		char, _ := types.ProtoResultToChar(result) // TODO: handle error

		w, _ := c.NextWriter(websocket.TextMessage) // TODO: handle error
		components.OOBWrap("beforeend:#submission-results", components.Text(char)).Render(ctx, w)
	}

	fmt.Println("close channel")
	c.WriteMessage(websocket.CloseMessage, nil)
}
