@echo off
REM ══════════════════════════════════════════════════════════
REM  AI Learning Platform — Development Server
REM  Usage: scripts\dev.bat
REM  Requirements: Docker Desktop
REM ══════════════════════════════════════════════════════════

cd /d "%~dp0\.."

REM Create .env from example if it doesn't exist
if not exist .env (
    echo Creating .env from .env.example...
    copy .env.example .env
    echo .env created. Edit it with your secrets if needed.
)

echo.
echo Starting AI Learning Platform (dev mode)...
echo    Backend:  http://localhost:8080
echo    Frontend: http://localhost:3000
echo.
echo    Hot-reload is enabled for both backend and frontend.
echo    Press Ctrl+C to stop all services.
echo.

docker compose -f docker-compose.dev.yml up --build %*
