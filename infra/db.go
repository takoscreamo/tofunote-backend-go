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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func SetupDB() *gorm.DB {
	env := getEnvOrDefault("ENV", "dev")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "ginuser")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "ginpassword")
	dbName := getEnvOrDefault("DB_NAME", "emotra")
	dbPort := getEnvOrDefault("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
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
