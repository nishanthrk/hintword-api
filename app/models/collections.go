package models

import "time"

// Collections [...]
type Collections struct {
	ID        string    `gorm:"primaryKey;column:id" json:"-"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Users     Users     `gorm:"joinForeignKey:user_id;foreignKey:id;references:UserID" json:"usersList"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.获取数据库表名
func (m *Collections) TableName() string {
	return "collections"
}

// CollectionsColumns get sql column name.获取数据库列名
var CollectionsColumns = struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	UserID:    "user_id",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
