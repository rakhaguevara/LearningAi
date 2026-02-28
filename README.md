# Adaptive AI Learning Platform — NeuraLearn

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

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)
- Redis 7 (or use Docker)

### 1. Clone & Configure

```bash
cp .env.example .env
# Edit .env with your credentials
```

### 2. Start Infrastructure

```bash
make docker-infra
```

This starts PostgreSQL and Redis in Docker containers.

### 3. Run Backend

```bash
make dev-backend
```

The API server starts at `http://localhost:8080`. Migrations run automatically on startup.

### 4. Run Frontend

```bash
cd frontend && npm install
make dev-frontend
```

The frontend starts at `http://localhost:3000`.

### 5. Run Everything Together

```bash
make dev
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
