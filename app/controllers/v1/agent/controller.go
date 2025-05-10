package agent_controller

import (
	"github.com/gofiber/fiber/v2"
	"hintword.com/api/app/common/validator"
	"net/http"

	agentService "hintword.com/api/app/services/agent"
)

func Completion(c *fiber.Ctx) error {
	params := ChatPayload
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	response, err := agentService.Execute(params.Context, params.Note)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"data":   response,
	})
}
