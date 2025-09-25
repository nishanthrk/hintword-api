package models

import (
	"time"

	"gorm.io/gorm"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
)

// Collections [...]
type Collections struct {
	CollectionID string    `gorm:"primaryKey;column:collection_id" json:"collection_id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	Users        Users     `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"-"`
	Name         string    `gorm:"column:name" json:"name"`
	Status       string    `gorm:"column:status" json:"status"`
	Sequence     int64     `gorm:"column:sequence" json:"sequence"`
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
	Sequence     string
	CreatedAt    string
	UpdatedAt    string
}{
	CollectionID: "collection_id",
	UserID:       "user_id",
	Name:         "name",
	Status:       "status",
	Sequence:     "sequence",
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

	// Set sequence to current timestamp for new collections
	if m.Sequence == 0 {
		m.Sequence = time.Now().Unix()
	}
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

func (m *Collections) FindByUserOrderedBySequence(userId string) (result []Collections, err error) {
	err = database.MysqlDB.Model(m).
		Where("`user_id` = ?", userId).
		Where("`status` = ?", StatusActive).
		Order("`sequence` ASC").
		Find(&result).Error
	return
}

func (m *Collections) UpdateSequence(collectionId string, sequence int64) (err error) {
	err = database.MysqlDB.Model(m).
		Where("`collection_id` = ?", collectionId).
		Update("sequence", sequence).Error
	return
}

func (m *Collections) BulkUpdateSequences(updates []struct {
	CollectionID string `json:"collection_id"`
	Sequence     int64  `json:"sequence"`
}) (err error) {
	// Use transaction for bulk update
	tx := database.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, update := range updates {
		if err := tx.Model(m).
			Where("`collection_id` = ?", update.CollectionID).
			Update("sequence", update.Sequence).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
