package models

import (
	"gorm.io/gorm"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"time"
)

// Tabs [...]
type Tabs struct {
	TabID        string      `gorm:"primaryKey;column:tab_id" json:"tab_id"`
	CollectionID string      `gorm:"column:collection_id" json:"collection_id"`
	Collections  Collections `gorm:"joinForeignKey:collection_id;foreignKey:collection_id;references:CollectionID" json:"-"`
	UserID       string      `gorm:"column:user_id" json:"user_id,omitempty"`
	Users        Users       `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"-"`
	Title        string      `gorm:"column:title" json:"title"`
	URL          string      `gorm:"column:url" json:"url"`
	FaviconURL   string      `gorm:"column:favicon_url" json:"favicon_url"`
	Sequence     int64       `gorm:"column:sequence" json:"sequence"`
	Status       string      `gorm:"column:status" json:"status"`
	CreatedAt    time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName get sql table name.
func (m *Tabs) TableName() string {
	return "tabs"
}

// TabsColumns get sql column name.
var TabsColumns = struct {
	TabID        string
	CollectionID string
	UserID       string
	Title        string
	URL          string
	FaviconURL   string
	Sequence     string
	Status       string
	CreatedAt    string
	UpdatedAt    string
}{
	TabID:        "tab_id",
	CollectionID: "collection_id",
	UserID:       "user_id",
	Title:        "title",
	URL:          "url",
	FaviconURL:   "favicon_url",
	Sequence:     "sequence",
	Status:       "status",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

func (m *Tabs) BeforeCreate(tx *gorm.DB) (err error) {
	if m.TabID == "" {
		m.TabID = utility.GenerateUUID()
	}
	m.CreatedAt = time.Now()
	m.Status = StatusActive
	m.UpdatedAt = time.Now()
	return
}

func (m *Tabs) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func (m *Tabs) Save() (result Tabs, err error) {
	err = database.MysqlDB.Save(m).Error
	return *m, err
}

func (m *Tabs) Create() (result Tabs, err error) {
	err = database.MysqlDB.Create(&m).Error
	return
}

func (m *Tabs) FindById(tabID string) (result Tabs, err error) {
	err = database.MysqlDB.Model(m).Where("`tab_id` = ?", tabID).Find(&result).Error
	return
}

func (m *Tabs) FindByCollectionId(collectionId string) (result []Tabs, err error) {
	err = database.MysqlDB.Model(m).
		Where("`collection_id` = ?", collectionId).
		Order("`sequence` asc").
		Find(&result).Error
	return
}
