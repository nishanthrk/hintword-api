package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"hintword.com/api/app/controllers/v1/note"
	"hintword.com/api/app/middlewares"
)

func SetupWsRoute(app *fiber.App) {
	ws := app.Group("/ws", middlewares.RequireLoggedIn())

	ws.Get("/note-sync", websocket.New(note.HandleNoteSync))
}
