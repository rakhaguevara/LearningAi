# ============================================
# Adaptive AI Learning Platform - Dev Launcher
# For Windows (PowerShell)
# ============================================
# Usage: .\dev.ps1
# This script replaces 'make dev' for Windows users.

Write-Host "🚀 Starting infrastructure (postgres + redis)..." -ForegroundColor Cyan
docker compose up -d postgres redis

Write-Host "⏳ Waiting for services to be healthy..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host "✅ Infrastructure ready. Starting backend & frontend..." -ForegroundColor Green

# Start backend in a new terminal window
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd backend; go run ./cmd/server"

# Start frontend in a new terminal window
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd frontend; npm run dev"

Write-Host ""
Write-Host "🎉 Development servers starting!" -ForegroundColor Green
Write-Host "   Backend  → http://localhost:8080" -ForegroundColor White
Write-Host "   Frontend → http://localhost:3000" -ForegroundColor White
Write-Host ""
Write-Host "Each service runs in its own terminal window." -ForegroundColor DarkGray
Write-Host "Close the terminal windows to stop the services." -ForegroundColor DarkGray
