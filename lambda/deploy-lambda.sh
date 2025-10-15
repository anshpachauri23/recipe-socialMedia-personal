#!/bin/bash

# Lambda Deployment Script for Recipe Social Media
# This script builds and deploys Lambda functions to AWS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting Lambda deployment...${NC}"

# Function to build and deploy a Lambda function
deploy_function() {
    local function_name=$1
    local function_path=$2
    local zip_name=$3
    
    echo -e "${YELLOW}📦 Building $function_name...${NC}"
    
    cd "$function_path"
    
    # Clean previous builds
    rm -f *.zip main
    
    # Build the Go binary
    echo "Building Go binary for $function_name..."
    GOOS=linux GOARCH=amd64 go build -o main .
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to build $function_name${NC}"
        exit 1
    fi
    
    # Create deployment package
    echo "Creating deployment package..."
    zip "$zip_name" main
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to create zip for $function_name${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ $function_name built successfully${NC}"
    
    # Return to parent directory
    cd - > /dev/null
}

# Deploy authentication functions
echo -e "${YELLOW}🔐 Deploying authentication functions...${NC}"

# Deploy register function
deploy_function "register" "functions/auth/register" "register-function.zip"

# Deploy login function  
deploy_function "login" "functions/auth/login" "login-function.zip"

# Deploy other functions if they exist
if [ -d "functions/posts/feed" ]; then
    echo -e "${YELLOW}📝 Deploying feed function...${NC}"
    deploy_function "feed" "functions/posts/feed" "feed-function.zip"
fi

if [ -d "functions/create-post" ]; then
    echo -e "${YELLOW}📝 Deploying create-post function...${NC}"
    deploy_function "create-post" "functions/create-post" "create-post-function.zip"
fi

if [ -d "functions/get-posts" ]; then
    echo -e "${YELLOW}📝 Deploying get-posts function...${NC}"
    deploy_function "get-posts" "functions/get-posts" "get-posts-function.zip"
fi

echo -e "${GREEN}🎉 All Lambda functions built successfully!${NC}"
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. Upload the zip files to your AWS Lambda functions"
echo "2. Make sure your Lambda functions have the correct environment variables:"
echo "   - DB_HOST"
echo "   - DB_PORT" 
echo "   - DB_USER"
echo "   - DB_PASSWORD"
echo "   - DB_NAME"
echo "   - JWT_SECRET"
echo "3. Test your functions to ensure they work correctly"

echo -e "${GREEN}✨ Deployment preparation complete!${NC}"
