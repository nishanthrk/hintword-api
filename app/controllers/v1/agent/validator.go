package agent_controller

var ChatPayload struct {
	Context string `json:"context" validate:"required,oneof=grammar optimisation"`
	Note    string `json:"note" validate:"required"`
	NoteId  string `json:"note_id" validate:"required"`
}
