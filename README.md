# Adaptive AI Learning Platform — Learny

AI-powered learning platform that personalizes teaching based on user interests.
Built with Next.js 14, Golang (Gin), PostgreSQL, Redis, and Alibaba Cloud Qwen LLM.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Client (Browser)                           │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │   Next.js 14 Frontend   │
                    │   (SSR + Client SPA)    │
                    │   TailwindCSS + Framer  │
                    └────────────┬────────────┘
                                 │ REST API
                    ┌────────────▼────────────┐
                    │   Gin HTTP Router       │
                    │   ┌──────────────────┐  │
                    │   │  Middleware Stack │  │
                    │   │  - CORS          │  │
                    │   │  - JWT Auth      │  │
                    │   │  - Rate Limiter  │  │
                    │   │  - Logger        │  │
                    │   │  - Recovery      │  │
                    │   └──────────────────┘  │
                    │                         │
                    │   ┌─── Modules ──────┐  │
                    │   │  auth            │  │
                    │   │  user            │  │
                    │   │  ai              │  │
                    │   │  learning        │  │
                    │   └─────────────────-┘  │
                    └───┬──────────┬──────────┘
                        │          │
              ┌─────────▼──┐  ┌───▼──────────┐
              │ PostgreSQL  │  │    Redis      │
              │  (primary)  │  │   (cache)     │
              └─────────────┘  └──────────────┘
                        │
              ┌─────────▼─────────────┐
              │  Alibaba Cloud        │
              │  ┌─────────────────┐  │
              │  │  Qwen LLM API   │  │
              │  │  OSS Storage    │  │
              │  └─────────────────┘  │
              └───────────────────────┘
```

## Folder Structure

```
/
├── backend/
│   ├── cmd/server/          # Application entrypoint
│   ├── internal/
│   │   ├── auth/            # Google OAuth + JWT authentication
│   │   ├── user/            # User profile management
│   │   ├── ai/              # AI provider abstraction + Qwen implementation
│   │   ├── learning/        # Learning session management
│   │   ├── middleware/       # Auth, CORS, rate limiting, logging
│   │   ├── config/          # Environment configuration
│   │   ├── database/        # PostgreSQL connection + migrations
│   │   ├── models/          # Domain models
│   │   └── common/          # Shared errors, response, logger
│   ├── pkg/jwt/             # JWT token service
│   ├── api/routes/          # Route registration
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app/             # Next.js App Router pages
│   │   ├── components/
│   │   │   ├── ui/          # AnimatedButton, GradientCard, etc.
│   │   │   ├── layout/      # Navbar, Footer
│   │   │   └── sections/    # Hero, Features, HowItWorks, Demo, CTA
│   │   ├── lib/             # Constants, utilities
│   │   ├── styles/          # Global CSS + Tailwind
│   │   └── types/           # TypeScript type definitions
│   ├── Dockerfile
│   ├── package.json
│   └── tailwind.config.js
├── infrastructure/
│   ├── docker/
│   ├── nginx/
│   └── scripts/
├── docs/
├── docker-compose.yml
├── Makefile
├── .env.example
└── README.md
```

## API Endpoints

| Method | Path                         | Auth     | Description                    |
|--------|------------------------------|----------|--------------------------------|
| GET    | `/health`                    | Public   | Health check                   |
| POST   | `/auth/google`               | Public   | Google OAuth authentication    |
| GET    | `/user/profile`              | Required | Get user profile + interests   |
| POST   | `/ai/explain`                | Required | AI-powered concept explanation |
| POST   | `/ai/generate-illustration`  | Required | Generate contextual visuals    |
| POST   | `/learning/start-session`    | Required | Start a learning session       |

## Database Models

- **User** — Core user account (Google OAuth linked)
- **LearningProfile** — Learning preferences, difficulty level, goals
- **InterestTag** — User interests with weighted categories
- **LearningSession** — Tracked learning sessions with duration
- **AIInteractionHistory** — Full audit trail of AI interactions

## Getting Started

### Prerequisites

- **Docker Desktop** (WAJIB) — [Download di sini](https://www.docker.com/products/docker-desktop/)
- Go 1.22+ (hanya jika ingin develop tanpa Docker)
- Node.js 20+ (hanya jika ingin develop tanpa Docker)

### Quick Start (Cara Tercepat — Semua OS)

**1. Clone repository:**
```bash
git clone <repository-url>
cd alibabab
```

**2. Buat file `.env`:**

| OS | Command |
|---|---|
| **Mac/Linux** | `cp .env.example .env` |
| **Windows CMD** | `copy .env.example .env` |
| **Windows PowerShell** | `Copy-Item .env.example .env` |

> **PENTING:** File `.env` TIDAK ikut ter-clone karena ada di `.gitignore`.
> Anda WAJIB membuat file `.env` dari `.env.example` sebelum menjalankan apapun!

**3. Jalankan project:**

| OS | Command |
|---|---|
| **Mac/Linux** | `make dev` atau `./scripts/dev.sh` |
| **Windows PowerShell** | `.\dev.ps1` |
| **Windows CMD** | `scripts\dev.bat` |

> Semua command di atas menjalankan Docker Compose, jadi Anda **tidak perlu install Go atau Node.js**.

**4. Buka di browser:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

### Manual Mode (Tanpa Docker untuk Backend/Frontend)

Jika ingin develop tanpa Docker (perlu Go + Node.js lokal):

```bash
# 1. Start database & redis via Docker
make docker-infra

# 2. Ubah .env: ganti DB_HOST=localhost dan REDIS_HOST=localhost

# 3. Run backend
make dev-backend

# 4. Run frontend (terminal baru)
cd frontend && npm install
make dev-frontend
```

### Docker (Full Stack)

```bash
make docker-up
```


## Available Make Commands

| Command           | Description                             |
|-------------------|-----------------------------------------|
| `make dev`        | Run backend + frontend concurrently     |
| `make build`      | Build both backend and frontend         |
| `make test`       | Run all tests                           |
| `make lint`       | Lint all code                           |
| `make docker-up`  | Start all services in Docker            |
| `make docker-down`| Stop all Docker services                |
| `make docker-infra`| Start only PostgreSQL + Redis          |
| `make clean`      | Remove build artifacts                  |

## Development Roadmap

### Phase 1 — Foundation (Current)
- [x] Project architecture and scaffolding
- [x] Backend API with auth, user, AI, and learning modules
- [x] Frontend landing page with animated sections
- [x] Docker infrastructure setup
- [ ] Connect Qwen LLM API (replace mock provider)
- [ ] Implement Google OAuth flow end-to-end
- [ ] Add frontend authentication state management

### Phase 2 — Core Product
- [ ] Interactive learning session UI
- [ ] Real-time AI chat interface
- [ ] Interest onboarding flow
- [ ] AI illustration generation via Qwen VL
- [ ] Learning progress dashboard
- [ ] Session history and replay

### Phase 3 — Intelligence Layer
- [ ] Adaptive difficulty adjustment
- [ ] Learning style detection from interaction patterns
- [ ] Spaced repetition integration
- [ ] Multi-modal content (video generation)
- [ ] Collaborative learning sessions

### Phase 4 — Scale
- [ ] Microservices extraction (AI service, session service)
- [ ] Event-driven architecture with message queues
- [ ] CDN integration for generated assets
- [ ] Multi-language support
- [ ] Mobile application (React Native)
- [ ] Enterprise/classroom features

## Tech Stack

| Layer          | Technology                                |
|----------------|------------------------------------------|
| Frontend       | Next.js 14, TypeScript, TailwindCSS, Framer Motion |
| Backend        | Go 1.22, Gin, PostgreSQL, Redis          |
| AI             | Alibaba Cloud Qwen LLM                  |
| Storage        | Alibaba OSS                              |
| Auth           | Google OAuth 2.0, JWT                    |
| Infrastructure | Docker, Docker Compose                   |
