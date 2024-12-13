package tasks

import (
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/DeepAung/gradient/public-server/views/components"
	"github.com/gofiber/fiber/v2"
)

var ErrInvalidPage = fiber.NewError(fiber.StatusBadRequest, "invalid page number")

const itemsPerPage = 50

type tasksHandler struct {
	tasksSvc types.TasksSvc
}

func InitTasksHandler(router fiber.Router, mid types.Middleware, tasksSvc types.TasksSvc) {
	handler := &tasksHandler{
		tasksSvc: tasksSvc,
	}

	tasksGroup := router.Group("/tasks")
	tasksGroup.Post("/", mid.OnlyAuthorized(), handler.GetTasks)
}

func (h *tasksHandler) GetTasks(c *fiber.Ctx) error {
	payload, ok := utils.GetPayload(c)
	if !ok {
		utils.DeleteTokenCookies(c)
		c.Response().Header.Add("HX-Redirect", "/signin")
		return nil
	}

	var dto types.GetTasksDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.RenderAlert(c, components.AlertError(err.Error()))
	}
	if err := utils.Validate(&dto); err != nil {
		return utils.RenderAlert(c, components.AlertError(err.Error()))
	}
	if dto.Page < 1 {
		return utils.RenderAlert(c, components.AlertError(ErrInvalidPage.Error()))
	}

	startIndex := itemsPerPage * (dto.Page - 1)
	stopIndex := startIndex + itemsPerPage
	tasks, err := h.tasksSvc.GetTasks(
		payload.UserId,
		dto.Search,
		dto.OnlyCompleted,
		startIndex,
		stopIndex,
	)
	if err != nil {
		_, msg := utils.ParseError(err)
		return utils.RenderAlert(c, components.AlertError(msg))
	}

	if len(tasks) == 0 {
		c.Response().Header.Add("HX-Trigger", "setMaxPage")
		return utils.Render(c, components.TaskNotFound())
	}

	return utils.Render(c, components.TasksTable(tasks))
}
