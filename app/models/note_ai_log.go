package models

import (
	"github.com/guregu/null"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"time"
)

// NoteAiLogs [...]
type NoteAiLogs struct {
	LogID           string            `gorm:"primaryKey;column:log_id" json:"-"`
	NoteID          string            `gorm:"column:note_id" json:"noteId"`
	Notes           Notes             `gorm:"joinForeignKey:note_id;foreignKey:note_id;references:NoteID" json:"notesList"`
	UserID          string            `gorm:"column:user_id" json:"userId"`
	Users           Users             `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"usersList"`
	InteractionType string            `gorm:"column:interaction_type" json:"interactionType"`
	RequestPayload  string            `gorm:"column:request_payload" json:"requestPayload"`
	ResponsePayload string            `gorm:"column:response_payload" json:"responsePayload"`
	TokenUsage      datatypes.JSONMap `gorm:"column:token_usage" json:"tokenUsage"`
	AudioURL        null.String       `gorm:"column:audio_url" json:"audioUrl"`
	ModelUsed       null.String       `gorm:"column:model_used" json:"modelUsed"`
	LanguageCode    null.String       `gorm:"column:language_code" json:"languageCode"`
	CreatedAt       time.Time         `gorm:"column:created_at" json:"createdAt"`
}

// TableName get sql table name.
func (m *NoteAiLogs) TableName() string {
	return "note_ai_logs"
}

// NoteAiLogsColumns get sql column name.
var NoteAiLogsColumns = struct {
	LogID           string
	NoteID          string
	UserID          string
	InteractionType string
	RequestPayload  string
	ResponsePayload string
	TokenUsage      string
	AudioURL        string
	ModelUsed       string
	LanguageCode    string
	CreatedAt       string
}{
	LogID:           "log_id",
	NoteID:          "note_id",
	UserID:          "user_id",
	InteractionType: "interaction_type",
	RequestPayload:  "request_payload",
	ResponsePayload: "response_payload",
	TokenUsage:      "token_usage",
	AudioURL:        "audio_url",
	ModelUsed:       "model_used",
	LanguageCode:    "language_code",
	CreatedAt:       "created_at",
}

func (m *NoteAiLogs) BeforeCreate(tx *gorm.DB) (err error) {
	if m.LogID == "" {
		m.LogID = utility.GenerateUUID()
	}
	m.CreatedAt = time.Now()
	return
}

func (m *NoteAiLogs) BeforeUpdate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	return
}

func (m *NoteAiLogs) Save() (result NoteAiLogs, err error) {
	err = database.MysqlDB.Save(m).Error
	return *m, err
}
