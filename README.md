# EmoTra Backend
- メンタルヘルスのための感情トラッキングアプリケーション
- 日記のように感情を記録し、振り返ることができる
- GoとGinによるバックエンドAPI

## 環境構築
- ローカルにGo、golang-migrateのインストール
- リポジトリをクローン
- リポジトリのディレクトリに移動
- `cp .env.example .env` で環境変数ファイルを作成
- `docker compose up -d --build` でコンテナを立ち上げ
- ブラウザで`http://localhost:8080/ping`にアクセス
- ビルド`go run main.go`

## ホットリロードで立ち上げ
- `go install github.com/air-verse/air@latest`でairをインストールしていること
- `air init`で設定ファイルを生成
- `air`で起動

## DBマイグレーションのupコマンド
- `migrate -database "mysql://root:password@tcp(localhost:3306)/emo_tracking" -path db/migrations up`
- `migrate -database "mysql://root:password@tcp(localhost:3306)/emo_tracking" -path db/migrations down`

## アーキテクチャなどの構想
- オニオンアーキテクチャを採用する
- ドメイン駆動設計を意識したい
- テストを書きやすくしたい
  - ユースケース層とリポジトリ層はインターフェース経由で依存関係を注入
  - コントローラ層はインターフェースを作成しないことにしてみる


## GoとGinの初期構築メモ
- `make emotra-backend`
- `cd emotra-backend`
- `go mod init emotra-backend`
- `go get -u github.com/gin-gonic/gin`
- `main.go`を作成
- `go run main.go`で起動
- `http://localhost:8080/ping`にアクセスして確認
- `{"message":"pong"}`が返ってくれば成功
