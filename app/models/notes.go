package models

import (
	"gorm.io/datatypes"
	"time"
)

// Notes [...]
type Notes struct {
	ID        string         `gorm:"primaryKey;column:id" json:"-"`
	UserID    string         `gorm:"column:user_id" json:"userId"`
	Users     Users          `gorm:"joinForeignKey:user_id;foreignKey:id;references:UserID" json:"usersList"`
	Title     string         `gorm:"column:title" json:"title"`
	Content   string         `gorm:"column:content" json:"content"`
	Folder    string         `gorm:"column:folder" json:"folder"`
	Tags      datatypes.JSON `gorm:"column:tags" json:"tags"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.获取数据库表名
func (m *Notes) TableName() string {
	return "notes"
}

// NotesColumns get sql column name.获取数据库列名
var NotesColumns = struct {
	ID        string
	UserID    string
	Title     string
	Content   string
	Folder    string
	Tags      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	UserID:    "user_id",
	Title:     "title",
	Content:   "content",
	Folder:    "folder",
	Tags:      "tags",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
