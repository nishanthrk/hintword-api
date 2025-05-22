package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	noteController "hintword.com/api/app/controllers/v1/note"
	"hintword.com/api/app/handler/agant"
	"hintword.com/api/app/middlewares"
)

func SetupWsRoute(app *fiber.App) {
	ws := app.Group("/ws", middlewares.RequireLoggedIn())

	ws.Get("/note-sync", websocket.New(noteController.HandleNoteSync))

	// Transcription WebSocket endpoint
	ws.Use("/transcribe", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	ws.Get("/transcribe", websocket.New(agant.HandleTranscriptionWebSocket))
}
