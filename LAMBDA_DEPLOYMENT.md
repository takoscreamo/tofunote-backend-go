# AWS Lambda デプロイガイド

このガイドでは、Emotra Backend APIをAWS Lambdaにデプロイする手順を説明します。

## 前提条件

1. **AWS CLI** がインストールされていること
2. **AWS SAM CLI** がインストールされていること
3. **Go 1.24.1以上** がインストールされていること
4. **AWS認証情報** が設定されていること

## セットアップ

### 1. AWS CLIの設定

```bash
aws configure
```

### 2. AWS SAM CLIのインストール

```bash
# macOS
brew install aws-sam-cli

# その他のOS
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
```

## デプロイ手順

### 1. 環境変数の設定

データベース接続情報を環境変数として設定します：

```bash
export DB_HOST="your-database-host"
export DB_PORT="5432"
export DB_USER="your-database-user"
export DB_PASSWORD="your-database-password"
export DB_NAME="emotra"
```

### 2. ビルドとデプロイ

```bash
# デプロイスクリプトを実行
./deploy.sh
```

または、手動で実行する場合：

```bash
# Lambda用のバイナリをビルド
./build-lambda.sh

# SAMでデプロイ
sam deploy \
  --template-file template.yaml \
  --stack-name emotra-api \
  --capabilities CAPABILITY_IAM \
  --region ap-northeast-1 \
  --parameter-overrides \
    DatabaseHost=$DB_HOST \
    DatabasePort=$DB_PORT \
    DatabaseUser=$DB_USER \
    DatabasePassword=$DB_PASSWORD \
    DatabaseName=$DB_NAME
```

### 3. デプロイ確認

デプロイが完了すると、API GatewayのエンドポイントURLが表示されます。

## ローカル開発

ローカル環境で開発する場合は、元のmain.goを使用します：

```bash
# ローカル用のmain.goにリネーム
mv main.go main_lambda.go
mv main_local.go main.go

# ローカルサーバーを起動
go run main.go
```

## API エンドポイント

デプロイ後、以下のエンドポイントが利用可能になります：

- `GET /api/me/diaries` - 日記一覧取得
- `GET /api/me/diaries/range` - 日付範囲で日記取得
- `GET /api/me/diaries/{date}` - 特定日付の日記取得
- `POST /api/me/diaries` - 日記作成
- `PUT /api/me/diaries/{date}` - 日記更新
- `DELETE /api/me/diaries/{date}` - 日記削除
- `GET /api/me/analyze-diaries` - 日記分析

## 注意事項

1. **データベース接続**: Lambda関数は外部データベース（RDS等）に接続する必要があります
2. **コールドスタート**: 初回実行時はコールドスタートが発生する可能性があります
3. **タイムアウト**: デフォルトで30秒のタイムアウトが設定されています
4. **メモリ**: デフォルトで512MBのメモリが割り当てられています

## トラブルシューティング

### ビルドエラー

```bash
# 依存関係を更新
go mod tidy

# キャッシュをクリア
go clean -cache
```

### デプロイエラー

```bash
# SAMのログを確認
sam logs -n EmotraAPIFunction --stack-name emotra-api --region ap-northeast-1

# CloudFormationスタックの状態を確認
aws cloudformation describe-stacks --stack-name emotra-api --region ap-northeast-1
```

### データベース接続エラー

1. セキュリティグループの設定を確認
2. VPC設定を確認
3. データベースの接続情報を確認 