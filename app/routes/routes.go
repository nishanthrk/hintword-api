package routes

import (
	"github.com/gofiber/fiber/v2"
	agentController "hintword.com/api/app/controllers/v1/agent"
	authController "hintword.com/api/app/controllers/v1/auth"
	noteController "hintword.com/api/app/controllers/v1/note"
	tabController "hintword.com/api/app/controllers/v1/tab"
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
	auth.Post("/google/user-info", authController.GoogleUserInfo)

	note := v1.Group("/note", middlewares.RequireLoggedIn())

	// Note routes
	note.Get("/list", noteController.GetNoteList)
	note.Post("/create", noteController.CreateUpdateNote)

	tab := v1.Group("/tab", middlewares.RequireLoggedIn())
	tab.Post("/collection", tabController.CreateUpdateCollection)
	tab.Get("/collection", tabController.GetCollection)
	tab.Post("/create", tabController.CreateUpdateTab)

	agent := v1.Group("/agent", middlewares.RequireLoggedIn())
	agent.Post("/completion", agentController.Completion)
}
