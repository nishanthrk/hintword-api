package models

import (
	"gorm.io/gorm"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/database"
	"time"
)

// Users [...]
type Users struct {
	ID        string    `gorm:"primaryKey;column:id" json:"-"`
	GoogleID  string    `gorm:"column:google_id" json:"googleId"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	AvatarURL string    `gorm:"column:avatar_url" json:"avatarUrl"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *Users) TableName() string {
	return "users"
}

// UsersColumns get sql column name.
var UsersColumns = struct {
	ID        string
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	GoogleID:  "google_id",
	Email:     "email",
	Name:      "name",
	AvatarURL: "avatar_url",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func (m *Users) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utility.GenerateUUID()
	}
	m.CreatedAt = time.Now()
	m.Status = StatusActive
	m.UpdatedAt = time.Now()
	return
}

func (m *Users) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func (m *Users) Save() (result Users, err error) {
	err = database.MysqlDB.Save(&m).Last(&result).Error
	return
}

func (m *Users) Create() (result Users, err error) {
	err = database.MysqlDB.Create(&m).Error
	return
}

func (m *Users) FindById(userId string) (result Users, err error) {
	err = database.MysqlDB.Model(m).Where("`id` = ?", userId).Find(&result).Error
	return
}

func (m *Users) FindByMobile(mobile string) (results Users, err error) {
	err = database.MysqlDB.Model(m).Where("`mobile` = ?", mobile).Find(&results).Error
	return
}

func (m *Users) FindByEmail(email string) (result Users, err error) {
	err = database.MysqlDB.Model(m).Where("`email` = ?", email).Find(&result).Error
	return
}
