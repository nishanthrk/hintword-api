package models

import "time"

// Tabs [...]
type Tabs struct {
	ID           string      `gorm:"primaryKey;column:id" json:"-"`
	CollectionID string      `gorm:"column:collection_id" json:"collectionId"`
	Collections  Collections `gorm:"joinForeignKey:collection_id;foreignKey:id;references:CollectionID" json:"collectionsList"`
	UserID       string      `gorm:"column:user_id" json:"userId"`
	Users        Users       `gorm:"joinForeignKey:user_id;foreignKey:id;references:UserID" json:"usersList"`
	Title        string      `gorm:"column:title" json:"title"`
	URL          string      `gorm:"column:url" json:"url"`
	FaviconURL   string      `gorm:"column:favicon_url" json:"faviconUrl"`
	CreatedAt    time.Time   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.获取数据库表名
func (m *Tabs) TableName() string {
	return "tabs"
}

// TabsColumns get sql column name.获取数据库列名
var TabsColumns = struct {
	ID           string
	CollectionID string
	UserID       string
	Title        string
	URL          string
	FaviconURL   string
	CreatedAt    string
	UpdatedAt    string
}{
	ID:           "id",
	CollectionID: "collection_id",
	UserID:       "user_id",
	Title:        "title",
	URL:          "url",
	FaviconURL:   "favicon_url",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}
