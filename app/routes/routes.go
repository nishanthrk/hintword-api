package routes

import (
	"github.com/gofiber/fiber/v2"
	authController "hintword.com/api/app/controllers/v1/auth"
	"hintword.com/api/app/middlewares"
)

func SetupRoutes(app *fiber.App) {
	welcomeFunction := func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to hintword API",
		})
	}

	welcome := app.Group("/welcome")

	welcome.Get("/", welcomeFunction)

	v1 := app.Group("/v1", middlewares.OptionalAuth())

	v1.Get("/", welcomeFunction)

	auth := v1.Group("/auth")
	auth.Get("/google", authController.GoogleAuth)
}
