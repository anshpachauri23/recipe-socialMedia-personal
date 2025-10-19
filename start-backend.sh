#!/bin/bash

# Start Backend Server Script
# This script starts the Go backend server with local environment configuration

echo "🚀 Starting RecipeShare Backend Server..."
echo "📍 Backend will run on: http://localhost:8080"
echo "🗄️  Database: AWS RDS PostgreSQL"
echo ""

# Navigate to backend directory
cd "$(dirname "$0")/backend"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check if local.env exists
if [ ! -f "local.env" ]; then
    echo "❌ Error: local.env file not found in backend directory."
    echo "Please create local.env with your database credentials."
    exit 1
fi

# Install dependencies if needed
echo "📦 Installing Go dependencies..."
go mod tidy

# Start the server
echo "🔄 Starting server..."
go run main.go
