# EmoTra Backend
- メンタルヘルスのための感情トラッキングアプリケーション
- 日記のように感情を記録し、振り返ることができる
- GoとGinによるバックエンドAPI

## 実装済み機能
- 日記の登録、編集、削除、閲覧
- グラフによる感情の可視化
- LLMにメンタルスコアを分析させる

## TODO
- JWT認証のログインを実装
- jsonbカラムを追加してスキーマレスにメンタル以外のスコアも記録できるようにする
- 設定系

## 環境構築

### ローカル開発環境

1. **前提条件**
   - Go 1.24.1以上
   - Docker & Docker Compose
   - golang-migrate

2. **セットアップ手順**
   ```bash
   # リポジトリをクローン
   git clone <repository-url>
   cd emotra-backend-go
   
   # 環境変数ファイルを作成
   cp .env.example .env
   
   # コンテナを立ち上げ
   docker compose up -d --build
   
   # ローカルサーバーを起動
   go run main_local.go
   ```

3. **動作確認**
   - ブラウザで`http://localhost:8080/ping`にアクセス
   - `{"message":"pong"}`が返ってくれば成功

### ホットリロードで立ち上げ
```bash
# airをインストール
go install github.com/air-verse/air@latest

# 設定ファイルを生成
air init

# ホットリロードで起動
air
```

## デプロイ

### AWS Lambda へのデプロイ

詳細な手順は [LAMBDA_DEPLOYMENT.md](./LAMBDA_DEPLOYMENT.md) を参照してください。

#### クイックデプロイ
```bash
# 環境変数を設定
export DB_HOST="your-database-host"
export DB_PORT="5432"
export DB_USER="your-database-user"
export DB_PASSWORD="your-database-password"
export DB_NAME="emotra"

# デプロイ実行
./deploy.sh
```

#### 前提条件
- AWS CLI
- AWS SAM CLI
- AWS認証情報の設定

## データベース

### マイグレーション
```bash
# マイグレーション実行（db接続情報は要変更）
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/emotra?sslmode=disable" up

# ロールバック
migrate -path ./infra/migrations -database "postgres://ginuser:ginpassword@localhost:5432/emotra?sslmode=disable" down
```

## テスト
```bash
# 全テスト実行
go test ./...
```

## API ドキュメント

### Swagger UI
- `http://localhost:8080/swagger` - Swagger UIで手動でAPIをテスト
- `http://localhost:8080/openapi.yml` - OpenAPI仕様書

### VSCode REST Client
- `.vscode/api-manual-test.http` を開く
- Send Request ボタンを押してAPIをテスト

## API エンドポイント

### 日記関連
- `GET /api/me/diaries` - 日記一覧取得
- `GET /api/me/diaries/range` - 日付範囲で日記取得
- `GET /api/me/diaries/{date}` - 特定日付の日記取得
- `POST /api/me/diaries` - 日記作成
- `PUT /api/me/diaries/{date}` - 日記更新
- `DELETE /api/me/diaries/{date}` - 日記削除

### 分析関連
- `GET /api/me/analyze-diaries` - 日記分析

## アーキテクチャ

### 設計原則
- **軽量DDD**を採用
- **オニオンアーキテクチャ**を採用
- **テストを書きやすく**設計

### 依存関係の注入
- ユースケース層とリポジトリ層はインターフェース経由で依存関係を注入
- コントローラ層はインターフェースを作成しない

## 開発メモ

### GoとGinの初期構築
```bash
make emotra-backend
cd emotra-backend
go mod init emotra-backend
go get -u github.com/gin-gonic/gin
# main.goを作成
go run main.go
```

### 便利なコマンド
```bash
# 8080ポートのプロセスを強制kill
kill -9 $(lsof -t -i:8080)

# データベースのテーブルを削除
# （必要に応じて実行）
```

## ファイル構成

```
emotra-backend-go/
├── api/controllers/          # コントローラー層
├── domain/                   # ドメイン層
├── infra/                    # インフラ層
├── repositories/             # リポジトリ層
├── routes/                   # ルーティング設定
├── usecases/                 # ユースケース層
├── main.go                   # Lambda用エントリーポイント
├── main_local.go             # ローカル開発用エントリーポイント
├── build-lambda.sh           # Lambdaビルドスクリプト
├── deploy.sh                 # デプロイスクリプト
├── template.yaml             # AWS SAMテンプレート
└── LAMBDA_DEPLOYMENT.md      # Lambdaデプロイガイド
```
