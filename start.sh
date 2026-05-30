#!/bin/bash

# Predit - Full Stack Startup Script
# This script starts the Backend, Frontend, and AI Referee services.

# Set default keys if not provided
export LIVEKIT_API_KEY=${LIVEKIT_API_KEY:-devkey}
export LIVEKIT_API_SECRET=${LIVEKIT_API_SECRET:-secretkey}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Predit Full Feature Stack...${NC}"

# Check for Docker
if command -v docker-compose >/dev/null 2>&1; then
    echo -e "${GREEN}🐳 Starting Infrastructure (Postgres, LiveKit)...${NC}"
    docker-compose up -d
else
    echo -e "${YELLOW}⚠️ docker-compose not found. Skipping infra (Postgres/LiveKit Server).${NC}"
fi

# Function to cleanup background processes on exit
cleanup() {
    echo -e "\n${YELLOW}🛑 Shutting down services...${NC}"
    kill $(jobs -p)
    if command -v docker-compose >/dev/null 2>&1; then
        docker-compose down
    fi
    exit
}

trap cleanup SIGINT SIGTERM

# 1. Start Go Backend
echo -e "${GREEN}📦 Starting Go Backend on :8080...${NC}"
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

# 2. Start AI Referee (Python)
echo -e "${GREEN}🤖 Starting AI Referee on :8000...${NC}"
source ai-referee/venv/bin/activate
python ai-referee/main.py &
AI_PID=$!
deactivate

# 3. Start Frontend (Vite)
echo -e "${GREEN}💻 Starting React Frontend on :5173...${NC}"
cd frontend
# Check if node_modules exists, if not install
if [ ! -d "node_modules" ]; then
    npm install
fi
npm run dev -- --clearScreen false &
FRONTEND_PID=$!
cd ..

echo -e "${BLUE}✨ All services are running!${NC}"
echo -e "   - Frontend: http://localhost:5173"
echo -e "   - Backend API: http://localhost:8080"
echo -e "   - AI Referee: http://localhost:8000"
echo -e "${YELLOW}Press Ctrl+C to stop all services.${NC}"

# Wait for background processes
wait
