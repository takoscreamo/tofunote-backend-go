.PHONY: migrate-diary check-data convert-to-json migrate-from-json migrate-from-json-prod check-data-prod run build test dev

# 環境変数ファイルを読み込み（存在する場合）
ifneq (,$(wildcard .env))
    include .env
    export
endif

# 環境変数
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= ginuser
DB_PASS ?= ginpassword
DB_NAME ?= tofunote
DB_SSL_MODE ?= disable

# データベース接続文字列
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# 日記データ移行（テキストから直接）
migrate-diary:
	go run scripts/migrate_diary_data.go

# テキストをJSONに変換
convert-to-json:
	go run scripts/convert_text_to_json.go

# JSONからDBに移行（SQLite）
migrate-from-json:
	go run scripts/migrate_json_to_db.go

# JSONからDBに移行（PostgreSQL）
migrate-from-json-prod:
	@echo "Using PostgreSQL connection:"
	@echo "  Host: $(DB_HOST)"
	@echo "  User: $(DB_USER)"
	@echo "  Database: $(DB_NAME)"
	@echo "  Port: $(DB_PORT)"
	ENV=prod DB_HOST=$(DB_HOST) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_PORT=$(DB_PORT) go run scripts/migrate_json_to_db.go

# 移行されたデータ確認（SQLite）
check-data:
	go run scripts/check_migrated_data.go

# 移行されたデータ確認（PostgreSQL）
check-data-prod:
	@echo "Using PostgreSQL connection:"
	@echo "  Host: $(DB_HOST)"
	@echo "  User: $(DB_USER)"
	@echo "  Database: $(DB_NAME)"
	@echo "  Port: $(DB_PORT)"
	ENV=prod DB_HOST=$(DB_HOST) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_PORT=$(DB_PORT) go run scripts/check_migrated_data.go

# アプリケーション実行
run:
	go run main.go

# ビルド
build:
	go build -o bin/tofunote-backend main.go

# テスト実行
test:
	go test ./...

# 依存関係の整理
tidy:
	go mod tidy

# コードフォーマット
fmt:
	go fmt ./...

# リント
lint:
	golangci-lint run

# ローカル開発サーバー起動
dev:
	ENV=dev go run cmd/local/main.go 