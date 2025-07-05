#!/bin/bash

# LocalStack initialization script
# This script creates the necessary AWS resources for development

set -e

echo "🚀 Initializing LocalStack resources..."

# Wait for LocalStack to be ready
echo "⏳ Waiting for LocalStack to be ready..."
until curl -s http://localhost:4566/health > /dev/null 2>&1; do
    sleep 2
done

echo "✅ LocalStack is ready!"

# Configure AWS CLI to use LocalStack
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566

echo "📦 Creating S3 buckets..."
aws s3 mb s3://videogrinder-uploads --endpoint-url=http://localhost:4566
aws s3 mb s3://videogrinder-outputs --endpoint-url=http://localhost:4566

echo "🗃️ Creating DynamoDB tables..."
aws dynamodb create-table \
    --table-name video-jobs \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=status,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --global-secondary-indexes \
        '[
            {
                "IndexName": "status-index",
                "KeySchema": [
                    {
                        "AttributeName": "status",
                        "KeyType": "HASH"
                    }
                ],
                "Projection": {
                    "ProjectionType": "ALL"
                },
                "ProvisionedThroughput": {
                    "ReadCapacityUnits": 5,
                    "WriteCapacityUnits": 5
                }
            }
        ]' \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url=http://localhost:4566

echo "📬 Creating SQS queues..."
aws sqs create-queue \
    --queue-name video-processing-queue \
    --endpoint-url=http://localhost:4566

aws sqs create-queue \
    --queue-name video-processing-dlq \
    --endpoint-url=http://localhost:4566

echo "🎉 LocalStack initialization completed!"
echo ""
echo "📋 Created resources:"
echo "   S3 Buckets: videogrinder-uploads, videogrinder-outputs"
echo "   DynamoDB Table: video-jobs"
echo "   SQS Queues: video-processing-queue, video-processing-dlq"
echo ""
echo "🔗 LocalStack endpoints:"
echo "   Main: http://localhost:4566"
echo "   Health: http://localhost:4566/health"
