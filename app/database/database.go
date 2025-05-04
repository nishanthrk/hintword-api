package database

import (
	"fmt"
	"log"

	"hintword.com/api/app/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() error {
	dsn := configs.GetMysqlConfig().GetMysqlConnectionInfo()

	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}
