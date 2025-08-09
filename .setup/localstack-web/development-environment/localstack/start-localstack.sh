#!/bin/bash

# LocalStack Sample Data Setup Script
# This script creates sample resources for testing the LocalStack Web interface

echo "🚀 Creating sample LocalStack resources..."

LOCALSTACK_ENDPOINT="http://localhost:4566"

# Create S3 buckets
echo "📦 Creating S3 buckets..."
awslocal s3 mb s3://sample-bucket-1
awslocal s3 mb s3://sample-bucket-2
awslocal s3 mb s3://images-bucket

# Upload sample files to S3
echo "📁 Uploading sample files..."
echo "Hello World!" | awslocal s3 cp - s3://sample-bucket-1/hello.txt
echo '{"message": "Sample JSON file"}' | awslocal s3 cp - s3://sample-bucket-1/sample.json

# Create SNS topics
echo "🔔 Creating SNS topics..."
awslocal sns create-topic --name create-order-topic
awslocal sns create-topic --name order-created-topic

# Create SQS queues
echo "📝 Creating SQS queues..."
awslocal sqs create-queue --queue-name sample-queue
awslocal sqs create-queue --queue-name notifications-queue

awslocal sqs create-queue --queue-name orders-queue
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:orders-queue \
  --protocol sqs \
  --notification-endpoint arn:aws:sqs:us-east-1:000000000000:create-order-topic

# Send sample messages to SQS
echo "💌 Sending sample messages to SQS..."
awslocal sqs send-message \
  --queue-url "$LOCALSTACK_ENDPOINT/000000000000/sample-queue" \
  --message-body "Hello from SQS!"

awslocal sqs send-message \
  --queue-url "$LOCALSTACK_ENDPOINT/000000000000/orders-queue" \
  --message-body '{"orderId": "12345", "amount": 99.99, "customerId": "user-456"}'

# Create DynamoDB tables
echo "🗃️ Creating DynamoDB tables..."
awslocal dynamodb create-table \
  --table-name Users \
  --attribute-definitions AttributeName=userId,AttributeType=S \
  --key-schema AttributeName=userId,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST

awslocal dynamodb create-table \
  --table-name Orders \
  --attribute-definitions \
    AttributeName=orderId,AttributeType=S \
    AttributeName=customerId,AttributeType=S \
  --key-schema \
    AttributeName=orderId,KeyType=HASH \
    AttributeName=customerId,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST

# Add sample items to DynamoDB
echo "📊 Adding sample data to DynamoDB..."
awslocal dynamodb put-item \
  --table-name Users \
  --item '{
    "userId": {"S": "user-123"},
    "name": {"S": "João Silva"},
    "email": {"S": "joao@example.com"},
    "age": {"N": "30"}
  }'

awslocal dynamodb put-item \
  --table-name Users \
  --item '{
    "userId": {"S": "user-456"},
    "name": {"S": "Maria Santos"},
    "email": {"S": "maria@example.com"},
    "age": {"N": "25"}
  }'

awslocal dynamodb put-item \
  --table-name Orders \
  --item '{
    "orderId": {"S": "order-001"},
    "customerId": {"S": "user-123"},
    "amount": {"N": "99.99"},
    "status": {"S": "completed"},
    "createdAt": {"S": "2025-01-01T10:00:00Z"}
  }'

# Create Kinesis streams
echo "🌊 Creating Kinesis streams..."
awslocal kinesis create-stream \
  --stream-name user-events \
  --shard-count 1

awslocal kinesis create-stream \
  --stream-name application-logs \
  --shard-count 2

# Create KMS key
echo "🔒 Creating KMS key..."
awslocal kms create-key --tags '[{"TagKey":"_custom_id_","TagValue":"17AE4E28-0F42-454E-8B41-7A2421263172"}]'

# Send sample records to Kinesis
sleep 5  # Wait for streams to be active
echo "📡 Sending sample records to Kinesis..."
awslocal kinesis put-record \
  --stream-name user-events \
  --partition-key "user-123" \
  --data '{"eventType": "login", "userId": "user-123", "timestamp": "2025-01-01T10:00:00Z"}'

awslocal kinesis put-record \
  --stream-name user-events \
  --partition-key "user-456" \
  --data '{"eventType": "purchase", "userId": "user-456", "amount": 29.99, "timestamp": "2025-01-01T11:00:00Z"}'

# Create a simple Lambda function (if zip is available)
if command -v zip >/dev/null; then
  echo "⚡ Creating sample Lambda function..."
  
  # Create temporary directory for Lambda
  mkdir -p /tmp/lambda-sample
  cat > /tmp/lambda-sample/index.js << 'EOF'
exports.handler = async (event) => {
  console.log('Event:', JSON.stringify(event, null, 2));
  
  const response = {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from LocalStack Lambda!',
      input: event,
      timestamp: new Date().toISOString()
    })
  };
  
  return response;
};
EOF
  
  # Create zip file
  cd /tmp/lambda-sample
  zip function.zip index.js
  
  # Create Lambda function
  awslocal lambda create-function \
    --function-name sample-function \
    --runtime nodejs18.x \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --handler index.handler \
    --zip-file fileb://function.zip
  
  # Clean up
  rm -rf /tmp/lambda-sample
  
  echo "✅ Lambda function created successfully!"
else
  echo "⚠️  Zip not available, skipping Lambda function creation"
fi

echo ""
echo "🎉 Sample data setup complete!"
echo ""
echo "📋 Created resources:"
echo "   DynamoDB Tables: Users, Orders"
echo "   Kinesis Streams: user-events, application-logs"
echo "   KMS Key: 17AE4E28-0F42-454E-8B41-7A2421263172"
echo "   Lambda Functions: sample-function (if zip available)"
echo "   S3 Buckets: sample-bucket-1, sample-bucket-2, images-bucket"
echo "   SNS Topics: create-order-topic, order-created-topic"
echo "   SQS Queues: sample-queue, orders-queue, notifications-queue"
echo ""
echo "🌐 Open the LocalStack Web interface at http://localhost:3000 to explore!"

echo "localstack started"
