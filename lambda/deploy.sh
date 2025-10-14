#!/bin/bash

# Recipe Social Media Lambda Deployment Script

echo "🚀 Deploying Recipe Social Media API to AWS..."

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "❌ AWS CLI is not installed. Please install it first."
    exit 1
fi

# Check if SAM CLI is installed
if ! command -v sam &> /dev/null; then
    echo "❌ SAM CLI is not installed. Please install it first."
    echo "Install with: brew install aws-sam-cli"
    exit 1
fi

# Set your AWS region
AWS_REGION=${AWS_REGION:-us-east-1}

echo "📍 Using AWS Region: $AWS_REGION"

# Build and deploy
echo "🔨 Building and deploying Lambda functions..."

sam build

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "📦 Deploying to AWS..."

sam deploy --guided

if [ $? -ne 0 ]; then
    echo "❌ Deployment failed"
    exit 1
fi

echo "✅ Deployment completed successfully!"
echo ""
echo "🔗 Your API endpoints are now available at:"
echo "   - Login: POST /auth/login"
echo "   - Register: POST /auth/register"
echo "   - Get Posts: GET /posts"
echo "   - Create Post: POST /posts"
echo "   - Feed: GET /posts/feed"
echo ""
echo "📝 Don't forget to:"
echo "   1. Update your frontend API_URL to point to the new API Gateway URL"
echo "   2. Set up your RDS Data API cluster"
echo "   3. Configure your database secrets in AWS Secrets Manager"
echo "   4. Run the database schema on your RDS instance"
