package postgres

import (
	"log"
	"os"

	"github.com/papijo/lms/internal/models/dbmodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize Database
func InitializeDB() (*gorm.DB, error) {
	//Get Environment Variables
	DB_USER := os.Getenv("DB_USER")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	// Construct the DSN using environment variables for PostgreSQL
	dsn := "host=" + DB_HOST + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " port=" + DB_PORT

	// Connect to the Database (PostgreSQL)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	log.Println("⚡️🚀 Conected to the database successfully.")

	return db, nil
}

// Close Database
func CloseDB(db *gorm.DB) error {
	//Close the Database connection when the application exits
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}

// Migration Schema
func MigrateSchema(db *gorm.DB) error {
	// Migrate the schema
	err := db.AutoMigrate(&dbmodels.User{}, &dbmodels.Biodata{})
	if err != nil {
		return err
	}

	return nil
}
