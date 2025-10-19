#!/bin/bash

# Start Frontend Server Script
# This script starts the Next.js frontend with local backend configuration

echo "🚀 Starting RecipeShare Frontend Server..."
echo "📍 Frontend will run on: http://localhost:3000"
echo "🔗 Backend API: http://localhost:8080/api"
echo ""

# Navigate to project root
cd "$(dirname "$0")"

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ Error: npm is not installed. Please install npm first."
    exit 1
fi

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing Node.js dependencies..."
    npm install
fi

# Set environment variable for local development
export NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Start the development server
echo "🔄 Starting Next.js development server..."
npm run dev
