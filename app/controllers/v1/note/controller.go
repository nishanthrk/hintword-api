package note_controller

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
	"hintword.com/api/app/common/validator"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"hintword.com/api/app/models"
	userService "hintword.com/api/app/services/user"
)

func CreateUpdateNote(c *fiber.Ctx) error {
	params := PayloadNote{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}
	userDetails := userService.GetUserObject(c)
	note := models.Notes{}
	note, _ = note.FindByUser(params.NoteID, userDetails.UserId)
	if note.NoteID == "" && params.NoteID != "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  "invalid note id given",
		})
	}

	if note.NoteID == "" {
		note.NoteID = utility.GenerateUUID()
		note.UserID = userDetails.UserId
	}
	note.Title = params.Title
	note.Content = params.Content
	note.Status = params.Status
	note.Sequence = null.IntFrom(params.Sequence)
	note, err := note.Save()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -3,
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"data":   note,
	})
}

// HandleNoteSync handles real-time note synchronization via WebSocket
func HandleNoteSync(c *websocket.Conn) {
	// Get user from context
	userDetails := userService.GetUserObject(c)

	defer func() {
		err := c.Close()
		if err != nil {
			return
		}
		fmt.Printf("WebSocket connection closed for user: %s\n", userDetails.UserId)
	}()

	for {
		// Read incoming message
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message for user %s: %v\n", userDetails.UserId, err)
			break
		}

		// Parse incoming message
		var noteMsg CreateMessage
		if err := json.Unmarshal(msg, &noteMsg); err != nil {
			fmt.Printf("Invalid message format for user %s: %v\n", userDetails.UserId, err)
			// Send error response back to client
			errorResponse := map[string]interface{}{
				"status": 0,
				"error":  "Invalid message format",
			}
			if responseMsg, err := json.Marshal(errorResponse); err == nil {
				err = c.WriteMessage(websocket.TextMessage, responseMsg)
				if err != nil {
					return
				}
			}
			continue
		}

		// Set user ID and timestamps
		noteMsg.Payload.UserID = userDetails.UserId
		now := time.Now()
		noteMsg.Payload.UpdatedAt = now

		var response map[string]interface{}

		// If noteId is present in payload, it's an update
		if noteMsg.Payload.NoteID != "" {
			// Update existing note
			if err := database.MysqlDB.Model(&models.Notes{}).
				Where("note_id = ? AND user_id = ?", noteMsg.Payload.NoteID, userDetails.UserId).
				Updates(&noteMsg.Payload).Error; err != nil {
				response = map[string]interface{}{
					"status": 0,
					"error":  fmt.Sprintf("Failed to update note: %v", err),
				}
			} else {
				response = map[string]interface{}{
					"status": 1,
					"note":   noteMsg.Payload,
				}
			}
		} else {
			// Create new note
			noteMsg.Payload.CreatedAt = now
			noteMsg.Payload.NoteID = utility.GenerateUUID()
			if err := database.MysqlDB.Create(&noteMsg.Payload).Error; err != nil {
				response = map[string]interface{}{
					"status": 0,
					"error":  fmt.Sprintf("Failed to create note: %v", err),
				}
			} else {
				response = map[string]interface{}{
					"status": 1,
					"note":   noteMsg.Payload,
				}
			}
		}

		// Send response back to client
		if responseMsg, err := json.Marshal(response); err == nil {
			if err := c.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
				fmt.Printf("Error sending response to user %s: %v\n", userDetails.UserId, err)
				break
			}
		}
	}
}

// GetNoteList handles getting a list of notes for the authenticated user
func GetNoteList(c *fiber.Ctx) error {
	userDetails := userService.GetUserObject(c)
	var notes []models.Notes

	if err := database.MysqlDB.Where("user_id = ?", userDetails.UserId).
		Find(&notes).
		Order("updated_at desc").
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": 0,
			"error":  fmt.Sprintf("Failed to fetch notes: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"notes":  notes,
	})
}
