package models

import (
	"time"

	"gorm.io/datatypes"
)

// Notes [...]
type Notes struct {
	NoteID    string         `gorm:"primaryKey;column:note_id" json:"note_id"`
	UserID    string         `gorm:"column:user_id" json:"user_id"`
	Users     Users          `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"users,omitempty"`
	Title     string         `gorm:"column:title" json:"title"`
	Content   string         `gorm:"column:content" json:"content"`
	Folder    string         `gorm:"column:folder" json:"folder"`
	Tags      datatypes.JSON `gorm:"column:tags" json:"tags"`
	Status    string         `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName get sql table name.
func (m *Notes) TableName() string {
	return "notes"
}

// NotesColumns get sql column name.
var NotesColumns = struct {
	ID        string
	UserID    string
	Title     string
	Content   string
	Folder    string
	Tags      string
	Status    string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	UserID:    "user_id",
	Title:     "title",
	Content:   "content",
	Folder:    "folder",
	Tags:      "tags",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
