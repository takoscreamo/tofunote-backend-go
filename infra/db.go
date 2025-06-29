package infra

import (
	"feelog-backend/infra/db"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	dbName := getEnvOrDefault("DB_NAME", "feelog")
	dbPort := getEnvOrDefault("DB_PORT", "5432")

	// 環境に応じてSSL設定を変更
	var sslMode string
	if env == "prod" {
		sslMode = "require"
	} else {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo connect_timeout=10 statement_timeout=30000 idle_in_transaction_session_timeout=30000",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
		sslMode,
	)

	log.Printf("[DEBUG] SetupDB: ENV=%s, DB_HOST=%s, DB_USER=%s, DB_NAME=%s, DB_PORT=%s", env, dbHost, dbUser, dbName, dbPort)
	log.Printf("[DEBUG] SetupDB: DSN=%s", dsn)

	var (
		database *gorm.DB
		err      error
	)

	log.Println("[DEBUG] SetupDB: PostgreSQLに接続します")
	database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("[ERROR] SetupDB: DB接続失敗: %v", err)
		panic("Failed to connect database")
	}

	// 接続プールの設定
	sqlDB, err := database.DB()
	if err != nil {
		log.Printf("[ERROR] SetupDB: DB取得失敗: %v", err)
		panic("Failed to get database")
	}

	// 接続プール設定
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("[DEBUG] SetupDB: PostgreSQL接続完了")

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
