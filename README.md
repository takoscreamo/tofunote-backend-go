# Feelog Backend

メンタルヘルスのための感情トラッキングアプリケーションのバックエンド（Go + Gin）

---

## 実装済み機能
- 日記の登録、編集、削除、閲覧
- 感情のグラフ可視化
- LLMによるメンタルスコア分析
- JWT認証によるユーザー管理

## 設計意図

- **JWT認証**: セッションレスでスケーラブルな認証を実現。APIサーバーはステートレスで、トークンには最小限の情報のみを格納。
- **ユーザー管理**: usersテーブルはメールアドレスをユニークキー、nicknameやOAuth拡張用カラムも設計段階から用意。パスワードはbcryptでハッシュ化。
- **API設計・テスト**: RESTful設計・OpenAPI/Swaggerで仕様を明示。テーブル駆動テストで正常系・異常系を網羅。
- **アーキテクチャ**: 軽量DDD・オニオンアーキテクチャを採用。GORMによるDBアクセス、リポジトリインターフェースで抽象化。依存関係注入でテスト容易性・保守性を向上。
- **運用・拡張性**: golang-migrateでマイグレーション管理。エラー時のトラブルシュート手順も明記。OAuthや属性追加も容易なスキーマ設計。

## アーキテクチャ
- **軽量DDD/オニオンアーキテクチャ**: ドメイン・ユースケース・リポジトリ・インフラ層を明確に分離。
- **依存関係注入**: ユースケース層とリポジトリ層はインターフェース経由で依存注入。
- **コントローラ層**: Ginのハンドラとして実装し、インターフェースは作成しない。

## 環境構築

### 前提条件
- Go 1.24.1以上
- Docker & Docker Compose
- golang-migrate

### セットアップ手順
```bash
# リポジトリをクローン
git clone <repository-url>
cd feelog-backend-go

# 環境変数ファイルを作成
cp .env.example .env

# コンテナを立ち上げ
docker compose up -d --build

# ローカルサーバーを起動
make dev
```

### 動作確認
- ブラウザで`http://localhost:8080/ping`にアクセスし、`{"message":"pong"}`が返ればOK

## 利用可能なコマンド

```bash
make dev        # ローカル開発サーバー起動
make run        # アプリケーション実行
make build      # ビルド
make test       # テスト実行
make tidy       # 依存関係の整理
make fmt        # コードフォーマット
make lint       # リント
make migrate-diary         # 日記データ移行（テキスト→DB）
make convert-to-json       # テキスト→JSON変換
make migrate-from-json     # JSON→DB移行（SQLite）
make migrate-from-json-prod # JSON→DB移行（PostgreSQL）
make check-data           # 移行データ確認（SQLite）
make check-data-prod      # 移行データ確認（PostgreSQL）
```

### ホットリロード
```bash
go install github.com/air-verse/air@latest
air init
air
```

## データベース・マイグレーション

```bash
# マイグレーション実行（db接続情報は要変更）
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/feelog?sslmode=disable" up
# ロールバック
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/feelog?sslmode=disable" down
```
- usersテーブル追加後は必ずマイグレーションを実行
- エラー時は`Dirty database version`等をforceコマンドで復旧

## API仕様・認証

### ユーザー登録・ログインAPI
- `POST /api/register` : ユーザー新規登録（メール・パスワード・ニックネーム）
  - リクエスト例:
    ```json
    {
      "email": "test@example.com",
      "password": "password123",
      "nickname": "テストユーザー"
    }
    ```
  - レスポンス例:
    ```json
    { "token": "<JWTトークン>" }
    ```
- `POST /api/login` : ログイン（メール・パスワード）
  - リクエスト例:
    ```json
    {
      "email": "test@example.com",
      "password": "password123"
    }
    ```
  - レスポンス例:
    ```json
    { "token": "<JWTトークン>" }
    ```

### JWT認証
- JWTトークンは`Authorization: Bearer <JWTトークン>`ヘッダで送信
- 有効期限や署名鍵は環境変数で管理

### Swagger UIでのAPI動作確認
- `http://localhost:8080/swagger` でAPI仕様・動作確認が可能

### 日記・分析API
- `GET /api/me/diaries` - 日記一覧取得
- `GET /api/me/diaries/range` - 日付範囲で日記取得
- `GET /api/me/diaries/{date}` - 特定日付の日記取得
- `POST /api/me/diaries` - 日記作成
- `PUT /api/me/diaries/{date}` - 日記更新
- `DELETE /api/me/diaries/{date}` - 日記削除
- `GET /api/me/analyze-diaries` - 日記分析

### OpenAPI/Swagger
- `http://localhost:8080/swagger` - Swagger UI
- `http://localhost:8080/openapi.yml` - OpenAPI仕様書
- `.vscode/api-manual-test.http` でVSCode REST Clientからもテスト可能

## テスト
- テーブル駆動テストで正常系・異常系を網羅
- `go test $(go list ./... | grep -v scripts)` で実行

## ファイル構成

```
feelog-backend-go/
├── api/controllers/          # コントローラー層
├── domain/                   # ドメイン層
├── infra/                    # インフラ層
├── repositories/             # リポジトリ層
├── routes/                   # ルーティング設定
├── usecases/                 # ユースケース層
├── scripts/                  # データ移行スクリプト
├── main.go                   # Lambda用エントリーポイント
├── cmd/local/main.go         # ローカル開発用エントリーポイント
├── build-lambda.sh           # Lambdaビルドスクリプト
├── deploy.sh                 # デプロイスクリプト
├── template.yaml             # AWS SAMテンプレート
├── Makefile                  # ビルド・実行コマンド
└── LAMBDA_DEPLOYMENT.md      # Lambdaデプロイガイド
```

## 開発メモ

- 8080ポートのプロセス強制kill: `kill -9 $(lsof -t -i:8080)`
- DBテーブル削除などは必要に応じて実行
- Gin導入: `go get -u github.com/gin-gonic/gin`


## TODO
- ロギング
- pgadminユーザ
- リフレッシュトークンに有効期限追加

---
