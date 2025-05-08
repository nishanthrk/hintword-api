package models

import (
	"github.com/guregu/null"
	"hintword.com/api/app/database"
	"time"

	"gorm.io/datatypes"
)

// Notes [...]
type Notes struct {
	NoteID    string          `gorm:"primaryKey;column:note_id" json:"note_id"`
	UserID    string          `gorm:"column:user_id" json:"user_id"`
	Users     Users           `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"-"`
	Title     string          `gorm:"column:title" json:"title"`
	Content   string          `gorm:"column:content" json:"content"`
	Folder    null.String     `gorm:"column:folder" json:"folder"`
	Tags      *datatypes.JSON `gorm:"column:tags" json:"tags"`
	Sequence  null.Int        `gorm:"column:sequence" json:"sequence"`
	Status    string          `gorm:"column:status" json:"status"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updated_at"`
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
	Sequence  string
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
	Sequence:  "sequence",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func (m *Notes) FindById(noteId string) (result Notes, err error) {
	err = database.MysqlDB.Model(m).Where("`note_id` = ?", noteId).Find(&result).Error
	return
}

func (m *Notes) FindByUser(noteId string, userId string) (result Notes, err error) {
	err = database.MysqlDB.Model(m).
		Where("`note_id` = ?", noteId).Where("`user_id` = ?", userId).
		Find(&result).Error
	return
}

func (m *Notes) Save() (result Notes, err error) {
	err = database.MysqlDB.Save(m).Error
	return *m, err
}
