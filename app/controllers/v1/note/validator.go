package note_controller

import "hintword.com/api/app/models"

type CreateMessage struct {
	NoteID  string       `json:"note_id,omitempty"`
	Payload models.Notes `json:"payload"`
}
