package auth_controller

import (
	"github.com/gofiber/fiber/v2"
	"hintword.com/api/app/common/validator"
	"net/http"

	userService "hintword.com/api/app/services/user"
)

func GoogleAuth(c *fiber.Ctx) error {
	params := GoogleAuthRequest
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	token, err := userService.HandleGoogleAuth(params.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"data": fiber.Map{
			"token": token,
		},
	})
}
