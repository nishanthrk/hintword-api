package models

import (
	"gorm.io/gorm"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"time"
)

// Collections [...]
type Collections struct {
	CollectionID string    `gorm:"primaryKey;column:collection_id" json:"collection_id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	Users        Users     `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"-"`
	Name         string    `gorm:"column:name" json:"name"`
	Status       string    `gorm:"column:status" json:"status"`
	Tabs         []Tabs    `gorm:"foreignKey:CollectionID;references:CollectionID" json:"tabs,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName get sql table name.
func (m *Collections) TableName() string {
	return "collections"
}

// CollectionsColumns get sql column name.
var CollectionsColumns = struct {
	CollectionID string
	UserID       string
	Name         string
	Status       string
	CreatedAt    string
	UpdatedAt    string
}{
	CollectionID: "collection_id",
	UserID:       "user_id",
	Name:         "name",
	Status:       "status",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

func (m *Collections) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CollectionID == "" {
		m.CollectionID = utility.GenerateUUID()
	}
	m.CreatedAt = time.Now()
	m.Status = StatusActive
	m.UpdatedAt = time.Now()
	return
}

func (m *Collections) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func (m *Collections) Save() (result Collections, err error) {
	err = database.MysqlDB.Save(m).Error
	return *m, err
}

func (m *Collections) Create() (result Collections, err error) {
	err = database.MysqlDB.Create(&m).Error
	return
}

func (m *Collections) FindById(collectionId string) (result Collections, err error) {
	err = database.MysqlDB.Model(m).Where("`collection_id` = ?", collectionId).Find(&result).Error
	return
}
