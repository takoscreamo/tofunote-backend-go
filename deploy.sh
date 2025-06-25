#!/bin/bash

# デプロイ設定
STACK_NAME="emotra-api"
REGION="ap-northeast-1"

echo "Building Lambda function..."
./build-lambda.sh

echo "Deploying to AWS..."
sam deploy \
  --template-file template.yaml \
  --stack-name $STACK_NAME \
  --capabilities CAPABILITY_IAM \
  --region $REGION \
  --parameter-overrides \
    DatabaseHost=$DB_HOST \
    DatabasePort=$DB_PORT \
    DatabaseUser=$DB_USER \
    DatabasePassword=$DB_PASSWORD \
    DatabaseName=$DB_NAME

echo "Deployment completed!"
echo "API Gateway URL:"
aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --region $REGION \
  --query 'Stacks[0].Outputs[?OutputKey==`EmotraAPI`].OutputValue' \
  --output text 