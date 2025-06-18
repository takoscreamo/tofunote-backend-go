package infra

import (
	"emotra-backend/infra/db"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	env := os.Getenv("ENV")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var (
		database *gorm.DB
		err      error
	)

	if env == "prod" {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		log.Println("Setup postgresql database")
	} else {
		database, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("Setup sqlite database")
	}
	if err != nil {
		panic("Failed to connect database")
	}

	// AutoMigrateでテーブルを作成
	err = database.AutoMigrate(&db.DiaryModel{})
	if err != nil {
		panic("Failed to migrate database")
	}

	return database
}
