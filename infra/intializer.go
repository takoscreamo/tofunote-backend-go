package infra

import (
	"github.com/joho/godotenv"
)

func Initialize() {
	// .envファイルがなくてもエラーを無視
	_ = godotenv.Load()
}
