# ============================================
# Adaptive AI Learning Platform - Dev Launcher
# For Windows (PowerShell)
# ============================================
# Usage: .\dev.ps1
# This script replaces 'make dev' for Windows users.
# Requirements: Docker Desktop

$ErrorActionPreference = "Stop"

# Navigate to project root (where this script lives)
$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ProjectRoot

# Create .env from example if it doesn't exist
if (-not (Test-Path ".env")) {
    Write-Host "📋 Creating .env from .env.example..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
    Write-Host "✅ .env created! Edit it with your secrets if needed." -ForegroundColor Green
    Write-Host ""
}

Write-Host ""
Write-Host "🚀 Starting AI Learning Platform (dev mode via Docker)..." -ForegroundColor Cyan
Write-Host "   Backend:  http://localhost:8080" -ForegroundColor White
Write-Host "   Frontend: http://localhost:3000" -ForegroundColor White
Write-Host ""
Write-Host "   Hot-reload is enabled for both backend and frontend." -ForegroundColor DarkGray
Write-Host "   Press Ctrl+C to stop all services." -ForegroundColor DarkGray
Write-Host ""

docker compose -f docker-compose.dev.yml up --build
