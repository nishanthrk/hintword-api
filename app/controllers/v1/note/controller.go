package note

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"hintword.com/api/app/models"
	"hintword.com/api/app/services/user"
)

type NoteMessage struct {
	NoteID  string       `json:"note_id,omitempty"`
	Payload models.Notes `json:"payload"`
}

// HandleNoteSync handles real-time note synchronization via WebSocket
func HandleNoteSync(c *websocket.Conn) {
	// Get user from context
	userDetails := user.GetUserObject(c)

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
		var noteMsg NoteMessage
		if err := json.Unmarshal(msg, &noteMsg); err != nil {
			fmt.Printf("Invalid message format for user %s: %v\n", userDetails.UserId, err)
			// Send error response back to client
			errorResponse := map[string]interface{}{
				"status": 0,
				"error":  "Invalid message format",
			}
			if responseMsg, err := json.Marshal(errorResponse); err == nil {
				c.WriteMessage(websocket.TextMessage, responseMsg)
			}
			continue
		}

		// Set user ID and timestamps
		noteMsg.Payload.UserID = userDetails.UserId
		now := time.Now()
		noteMsg.Payload.UpdatedAt = now

		var response map[string]interface{}

		// If noteId is present in payload, it's an update
		if noteMsg.Payload.ID != "" {
			// Update existing note
			if err := database.MysqlDB.Model(&models.Notes{}).
				Where("id = ? AND user_id = ?", noteMsg.Payload.ID, userDetails.UserId).
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
			noteMsg.Payload.ID = utility.GenerateUUID()
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
	userDetails := user.GetUserObject(c)
	var notes []models.Notes

	if err := database.MysqlDB.Where("user_id = ?", userDetails.UserId).Find(&notes).Error; err != nil {
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
