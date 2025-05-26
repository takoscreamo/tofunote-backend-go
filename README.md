# GoとGinの初期構築メモ
- `make emotra-backend`
- `cd emotra-backend`
- `go mod init emotra-backend`
- `go get -u github.com/gin-gonic/gin`
- `main.go`を作成
- `go run main.go`で起動
- `http://localhost:8080/ping`にアクセスして確認
- `{"message":"pong"}`が返ってくれば成功

# ホットリロード
- `go install github.com/air-verse/air@latest`でairをインストールしていること
- `air init`で設定ファイルを生成
- `air`で起動
