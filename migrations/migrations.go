package migrations

import (
	"log"

	"github.com/pressly/goose/v3"
	"hintword.com/api/app/database"
)

// RunMigrations runs all database migrations using Goose
func RunMigrations() error {
	// Get the database connection
	db, err := database.MysqlDB.DB()
	if err != nil {
		return err
	}

	if err = goose.SetDialect("mysql"); err != nil {
		log.Printf("Warning: Dialet: %v", err)
		return err
	}

	// Run Goose migrations
	if err = goose.Up(db, "./migrations/db"); err != nil {
		log.Printf("Warning: Failed to run Goose migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func RollbackMigrations() error {
	db, err := database.MysqlDB.DB()
	if err != nil {
		return err
	}

	if err := goose.Down(db, "./migrations/db"); err != nil {
		return err
	}

	return nil
}
