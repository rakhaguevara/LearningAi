#!/usr/bin/env bash
# ══════════════════════════════════════════════════════════
#  AI Learning Platform — Development Server
#  Usage: ./scripts/dev.sh
#  Requirements: Docker Desktop
# ══════════════════════════════════════════════════════════

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Create .env from example if it doesn't exist
if [ ! -f .env ]; then
    echo "📋 Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env created. Edit it with your secrets if needed."
fi

echo ""
echo "🚀 Starting AI Learning Platform (dev mode)..."
echo "   Backend:  http://localhost:8080"
echo "   Frontend: http://localhost:3000"
echo ""
echo "   Hot-reload is enabled for both backend and frontend."
echo "   Press Ctrl+C to stop all services."
echo ""

docker compose -f docker-compose.dev.yml up --build "$@"
