#!/bin/bash

# Start Development Environment Script
# This script starts both backend and frontend servers for local development

echo "🚀 Starting RecipeShare Development Environment..."
echo "📍 Frontend: http://localhost:3000"
echo "📍 Backend API: http://localhost:8080"
echo "🗄️  Database: AWS RDS PostgreSQL"
echo ""

# Function to cleanup background processes on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down servers..."
    jobs -p | xargs -r kill
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Start backend in background
echo "🔄 Starting backend server..."
cd "$(dirname "$0")/backend"
go run main.go &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 3

# Start frontend
echo "🔄 Starting frontend server..."
cd "$(dirname "$0")"
export NEXT_PUBLIC_API_URL=http://localhost:8080/api
npm run dev &
FRONTEND_PID=$!

# Wait for both processes
echo ""
echo "✅ Development environment started!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for any process to exit
wait
