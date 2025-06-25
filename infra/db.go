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
	log.Println("[DEBUG] SetupDB: 開始")
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

	log.Printf("[DEBUG] SetupDB: ENV=%s, DB_HOST=%s, DB_USER=%s, DB_NAME=%s, DB_PORT=%s", env, dbHost, dbUser, dbName, dbPort)
	log.Printf("[DEBUG] SetupDB: DSN=%s", dsn)

	var (
		database *gorm.DB
		err      error
	)

	if env == "prod" {
		log.Println("[DEBUG] SetupDB: PostgreSQLに接続します")
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		log.Println("[DEBUG] SetupDB: PostgreSQL接続完了")
	} else {
		log.Println("[DEBUG] SetupDB: SQLiteに接続します")
		database, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("[DEBUG] SetupDB: SQLite接続完了")
	}
	if err != nil {
		log.Printf("[ERROR] SetupDB: DB接続失敗: %v", err)
		panic("Failed to connect database")
	}

	log.Println("[DEBUG] SetupDB: AutoMigrate開始")
	// AutoMigrateでテーブルを作成
	err = database.AutoMigrate(&db.DiaryModel{})
	if err != nil {
		log.Printf("[ERROR] SetupDB: マイグレーション失敗: %v", err)
		panic("Failed to migrate database")
	}
	log.Println("[DEBUG] SetupDB: AutoMigrate完了")

	log.Println("[DEBUG] SetupDB: 正常終了")
	return database
}
