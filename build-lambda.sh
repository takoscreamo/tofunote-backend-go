#!/bin/bash

# AWS Lambda用のバイナリをビルド
echo "Building Lambda binary..."

# Linux用のバイナリをビルド（LambdaはLinux環境で動作）
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# バイナリをzipファイルにパッケージ
echo "Creating deployment package..."
zip function.zip bootstrap

echo "Build completed! function.zip is ready for deployment." 