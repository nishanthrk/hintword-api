package agent_controller

import (
	"github.com/gofiber/fiber/v2"
	"hintword.com/api/app/common/validator"
	"hintword.com/api/app/models"
	userService "hintword.com/api/app/services/user"
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

	userDetails := userService.GetUserObject(c)

	log := models.NoteAiLogs{}
	log.UserID = userDetails.UserId
	log.NoteID = params.NoteId
	log.InteractionType = models.InteractionTypeCompletion

	response, err := agentService.Completion(&log, params.Context, params.Note)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  err.Error(),
		})
	}

	log.TokenUsage = map[string]interface{}{
		"prompt_tokens":     response.Usage.PromptTokens,
		"completion_tokens": response.Usage.CompletionTokens,
		"total_tokens":      response.Usage.TotalTokens,
	}

	log, err = log.Save()

	return c.JSON(fiber.Map{
		"status": 1,
		"data": fiber.Map{
			"message": response.Choices[0].Message.Content,
			"log_id":  log.LogID,
		},
	})
}
