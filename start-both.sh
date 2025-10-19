#!/bin/bash

# Start Both Frontend and Backend Script
# This script starts both servers and ensures they're working

echo "🚀 Starting RecipeShare Development Environment..."
echo "📍 Frontend: http://localhost:3000"
echo "📍 Backend API: http://localhost:8080"
echo "🗄️  Database: AWS RDS PostgreSQL"
echo "📸 Image Storage: AWS S3"
echo ""

# Function to cleanup background processes on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down servers..."
    pkill -f "go run main.go" 2>/dev/null
    pkill -f "npm run dev" 2>/dev/null
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Start backend
echo "🔄 Starting backend server..."
cd "$(dirname "$0")/backend"
go run main.go &
BACKEND_PID=$!

# Wait for backend to start
echo "⏳ Waiting for backend to start..."
sleep 3

# Test backend
if curl -s "http://localhost:8080/api/users/search?q=test" > /dev/null; then
    echo "✅ Backend is running successfully"
else
    echo "❌ Backend failed to start"
    exit 1
fi

# Start frontend
echo "🔄 Starting frontend server..."
cd "$(dirname "$0")"
npm run dev &
FRONTEND_PID=$!

# Wait for frontend to start
echo "⏳ Waiting for frontend to start..."
sleep 5

# Test frontend
if curl -s "http://localhost:3000" > /dev/null; then
    echo "✅ Frontend is running successfully"
else
    echo "❌ Frontend failed to start"
    exit 1
fi

echo ""
echo "🎉 Development environment is ready!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for any process to exit
wait
