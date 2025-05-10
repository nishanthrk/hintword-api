package note_controller

import (
	"gorm.io/datatypes"
	"hintword.com/api/app/models"
)

type CreateMessage struct {
	NoteID  string       `json:"note_id,omitempty"`
	Payload models.Notes `json:"payload"`
}

type PayloadNote struct {
	NoteID   string         `json:"note_id,omitempty"`
	Title    string         `json:"title" validate:"required"`
	Content  datatypes.JSON `json:"content" validate:"required"`
	Sequence int64          `json:"sequence" validate:"required"`
	Status   string         `json:"status" validate:"oneof=ACTIVE INACTIVE"`
}
