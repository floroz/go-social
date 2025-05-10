# Tech Context

## Technologies Used

- **Backend:**
    - Go (version specified in `go.mod`)
    - Chi Router (for HTTP routing)
    - PostgreSQL (database)
- **Frontend:**
    - React
    - TypeScript
    - Vite (build tool and dev server)
    - Tailwind CSS (styling)
    - Shadcn/ui (UI components)
    - Zustand (state management)
    - React Query (data fetching and caching)
    - Vitest (unit/integration testing)
    - MSW (Mock Service Worker for API mocking in tests)
- **Database Migrations:**
    - golang-migrate (CLI tool)
- **Development Environment:**
    - Docker & Docker Compose (for running PostgreSQL and potentially other services)
    - Air (optional, for Go backend live reload)
    - Node.js (version specified in `frontend/package.json` or latest LTS)
    - npm (for frontend dependency management)
- **API Specification:**
    - OpenAPI 3
    - YAML

## Development Setup

- **Prerequisites:** Go, Node.js, npm, Docker, Docker Compose, Golang Migrate CLI, Air (optional).
- **Installation:**
    1. Clone repository.
    2. `go mod tidy` for backend dependencies.
    3. `cd frontend && npm install` for frontend dependencies.
- **Environment Variables:**
    - Backend: `.env.local` (from `.env.local.example`) for DB connection, JWT secret.
    - Frontend: `frontend/.env.development` and `frontend/.env.production` (from `frontend/.env.example`) for `VITE_API_BASE_URL`.
- **Database:** Start PostgreSQL using `docker compose up -d`.
- **Migrations:** Run using `make migrate-up`.
- **Running the Application:**
    - Concurrently (backend & frontend): `make dev`
    - Backend only (with Air): `air`
    - Backend only (standard): `make dev-be` or `go run ./cmd/main.go`
    - Frontend only: `make dev-fe` or `cd frontend && npm run dev`

## Technical Constraints

- Backend and frontend types must be kept in sync with the OpenAPI specification.
- Adherence to the defined API response structure (success with `data`, error with `errors`).
- Go tools managed via `go.mod` and `tools.go`.
- Node.js tools managed via `frontend/package.json`.

## Dependencies

- **Backend (Key Go Modules - see `go.mod` for full list):**
    - `github.com/go-chi/chi/v5`
    - `github.com/jackc/pgx/v5` (PostgreSQL driver)
    - `github.com/golang-jwt/jwt/v5`
    - `github.com/deepmap/oapi-codegen` (for code generation, via `tools.go`)
- **Frontend (Key NPM Packages - see `frontend/package.json` for full list):**
    - `react`, `react-dom`
    - `typescript`
    - `vite`
    - `tailwindcss`
    - `@radix-ui/*` (used by Shadcn/ui)
    - `zustand`
    - `@tanstack/react-query`
    - `vitest`, `msw`
    - `openapi-typescript` (for code generation)
    - `@redocly/cli` (for OpenAPI bundling)

## Tool Usage Patterns

- **`make`:** Centralized script for common development tasks (running dev servers, tests, migrations, code generation). See `Makefile`.
- **`go mod tidy`:** To manage Go dependencies.
- **`npm install` (in `frontend/`):** To manage frontend dependencies.
- **`docker compose`:** To manage services like PostgreSQL.
- **`golang-migrate`:** For database schema migrations.
- **`air`:** For live reloading of the Go backend during development.
- **`oapi-codegen`:** Generates Go types from `openapi-bundled.yaml`. Invoked via `make generate-types`.
- **`openapi-typescript`:** Generates TypeScript types from `openapi-bundled.yaml`. Invoked via `make generate-types`.
- **`@redocly/cli bundle`:** Bundles partial OpenAPI YAML files into `openapi-bundled.yaml`. Invoked via `make generate-types`.
- **Swagger UI:** Served by the backend at `/api/docs` for interactive API exploration.
