# TOFU NOTE Backend（Go + Gin）

豆腐メンタルの可視化・記録アプリケーション。日々のメンタルスコアや日記を記録し、グラフで振り返ることができる。本リポジトリはそのバックエンドAPI。

---

## 主な機能

- ユーザー登録・認証（JWT）
- 日記の登録・編集・削除・取得
- 日記データの範囲・日付指定取得
- 感情グラフ可視化用データ提供
- LLM（大規模言語モデル）による日記分析・メンタルスコア算出
- データ移行・マイグレーション（テキスト/JSON→DB）
- Swagger/OpenAPIによるAPI仕様公開
- テーブル駆動テストによる品質担保

---

## アーキテクチャ

- **軽量DDD/オニオンアーキテクチャ**  
  ドメイン・ユースケース・リポジトリ・インフラ層を明確に分離し、保守性・拡張性を重視。
- **依存関係注入**  
  ユースケース層とリポジトリ層はインターフェース経由でDI。
- **リポジトリ層**  
  domain層（例: `domain/user/repository.go` や `domain/diary/repository.go`）にリポジトリインターフェースのみを定義し、repositories層でGORMを用いた実装を提供。インフラ依存の実装はrepositories配下に集約。
- **DBアクセス**  
  GORMを利用し、リポジトリインターフェースで抽象化。
- **認証**  
  JWTによるセッションレス認証。パスワードはbcryptでハッシュ化。

---

## ディレクトリ構成

```
tofunote-backend-go/
├── api/controllers/          # コントローラー層（Ginハンドラ）
├── domain/                   # ドメイン層（ビジネスロジック・エンティティ、値オブジェクト、リポジトリIF）
│   ├── diary/                # 日記ドメイン（エンティティ・値オブジェクト・リポジトリIF）
│   │   ├── diary.go          # 日記エンティティ・値オブジェクトの定義、日記関連のドメインロジック
│   │   ├── mental.go         # メンタルスコア等の値オブジェクト・ロジック
│   │   ├── mental_test.go    # メンタル関連ロジックのテスト
│   │   └── repository.go     # 日記リポジトリのインターフェース定義
│   └── user/                 # ユーザードメイン（エンティティ・値オブジェクト・リポジトリIF）
│       ├── user.go           # ユーザーエンティティ・値オブジェクトの定義、ユーザー関連のドメインロジック
│       └── repository.go     # ユーザーリポジトリのインターフェース定義
├── infra/                    # インフラ層（DBアクセス・外部連携）
│   ├── db/                  # DB関連の実装（GORMモデル、DB操作ヘルパーなど）
│   │   ├── diary.go         # 日記データのDBモデル・操作
│   │   └── user.go          # ユーザーデータのDBモデル・操作
│   ├── diary_datas/         # 日記データの初期データや補助データ
│   ├── migrations/          # golang-migrate用のマイグレーションSQL
│   ├── db.go                # DB接続・初期化処理
│   ├── intializer.go        # 各種初期化処理（例：依存注入や設定ロード）
│   ├── jwt.go               # JWT生成・検証ロジック
│   └── jwt_test.go          # JWT関連のテスト
├── repositories/             # GORMを用いたリポジトリ実装（インターフェースはdomain層に定義）
│   ├── diary_repositories.go   # GORM依存の実装
│   └── user_repositories.go    # GORM依存の実装
├── routes/                   # ルーティング・ミドルウェア
│   ├── api.go               # APIルーティングの定義（エンドポイントとコントローラの紐付け）
│   ├── cors.go              # CORS（クロスオリジン）設定
│   ├── middleware/
│   │   ├── jwt_auth.go      # JWT認証ミドルウェア
│   │   └── jwt_auth_test.go # JWT認証ミドルウェアのテスト
│   └── swagger.go           # Swagger UIのルーティング・設定
├── usecases/                 # ユースケース層（アプリケーションサービス）
├── scripts/                  # データ移行・補助スクリプト
├── main.go                   # Lambda用エントリーポイント
├── cmd/local/main.go         # ローカル開発用エントリーポイント
├── infra/migrations/         # DBマイグレーションファイル
├── Makefile                  # ビルド・実行コマンド
├── openapi.yml               # OpenAPI仕様書
└── README.md                 # 本ファイル
```

---

## セットアップ

### 前提

- Go 1.24.1以上
- Docker & Docker Compose
- golang-migrate

### 初期セットアップ

```bash
git clone <repository-url>
cd tofunote-backend-go
cp .env.example .env
docker compose up -d --build
make dev
```

- ブラウザで `http://localhost:8080/ping` にアクセスし、`{"message":"pong"}` が返ればOK

---

## コマンド一覧

| コマンド                | 説明                                 |
|------------------------|--------------------------------------|
| make dev               | ローカル開発サーバー起動             |
| make run               | アプリケーション実行                 |
| make build             | ビルド                               |
| make test              | テスト実行                           |
| make tidy              | 依存関係の整理                       |
| make fmt               | コードフォーマット                   |
| make lint              | リント                               |
| make migrate-diary     | 日記データ移行（テキスト→DB）        |
| make convert-to-json   | テキスト→JSON変換                    |
| make migrate-from-json | JSON→DB移行（SQLite）                |
| make migrate-from-json-prod | JSON→DB移行（PostgreSQL）        |
| make check-data        | 移行データ確認（SQLite）              |
| make check-data-prod   | 移行データ確認（PostgreSQL）          |

---

## DBマイグレーション

```bash
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/tofunote?sslmode=disable" up
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/tofunote?sslmode=disable" down
```
- usersテーブル追加後は必ずマイグレーションを実行
- エラー時は `Dirty database version` などをforceコマンドで復旧

---


## OpenAPI/Swagger

- `http://localhost:8080/swagger` : Swagger UI
- `http://localhost:8080/openapi.yml` : OpenAPI仕様書

---

## テスト

- テーブル駆動テストで正常系・異常系を網羅
- `go test $(go list ./... | grep -v scripts)` で実行

---

## 開発メモ

- 8080ポートのプロセス強制kill: `kill -9 $(lsof -t -i:8080)`
- DBテーブル削除などは必要に応じて実行
- Gin導入: `go get -u github.com/gin-gonic/gin`
- ホットリロード: `go install github.com/air-verse/air@latest && air init && air`

---

## TODO

- ロギング強化
- pgadminユーザ管理
- リフレッシュトークンの有効期限追加
- 管理者画面

---
