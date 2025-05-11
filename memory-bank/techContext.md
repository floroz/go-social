# Technical Context

## Technologies Used

Based on the current file structure, the following technologies are likely in use:

### Backend
- **Language:** Go
- **Framework/Router:** (To be determined - `cmd/api/api.go` would contain this information. Common choices include Gin, Echo, or net/http)
- **Database:** PostgreSQL (inferred from `cmd/migrate/migrations` SQL files and `docker-compose.yaml` if it includes a Postgres service)
- **Migrations:** `golang-migrate/migrate` (inferred from migration file naming convention `00000X_*.up.sql` and `00000X_*.down.sql`)
- **Validation:** `go-playground/validator/v10` (used for struct validation, e.g., in DTOs within services)
- **API Specification:** OpenAPI (Swagger) - `openapi/openapi.yaml`

### Frontend
- **Language:** TypeScript
- **Framework/Library:** React (inferred from `App.tsx`, `main.tsx`, `vite.config.ts`)
- **Build Tool/Bundler:** Vite (inferred from `vite.config.ts`, `index.html` script type module)
- **UI Components:** Likely custom components, potentially with a library like Shadcn/UI (inferred from `frontend/components/ui/`)
- **Styling:** CSS (inferred from `index.css`), potentially Tailwind CSS if `tailwind.config.js` exists or is configured in `postcss.config.js`.
- **State Management:** (To be determined - `frontend/src/stores/authStore.ts` suggests a custom store or a library like Zustand or Jotai)
- **Testing:**
    - Unit/Integration: Vitest (inferred from `vitest.config.ts`)
    - E2E/Component: Playwright (inferred from `frontend/playwright/`)
- **API Client:** Custom fetch wrapper (`frontend/src/lib/api.ts`) or a library like Axios.
- **Mocking:** MSW (Mock Service Worker) (inferred from `frontend/src/mocks/`)

### General
- **Containerization:** Docker (`docker-compose.yaml`, `.air.toml` might be for live reloading with Docker)
- **Version Control:** Git (`.gitignore`)
- **Package Management:**
    - Go: Go Modules (`go.mod`, `go.sum`)
    - Node.js: npm or pnpm or yarn (`package.json`, `package-lock.json`)
- **Task Runner/Build System:** Make (`Makefile`)

## Development Setup

- **Backend:** Likely run using `go run cmd/main.go` or via Docker, possibly with live reload configured (e.g., Air via `.air.toml`).
- **Frontend:** Likely run using `npm run dev` (or `pnpm dev`/`yarn dev`) which would start the Vite development server.
- **Database:** Managed via Docker Compose. Migrations are applied manually or via a script.
- **Code Generation:**
    - OpenAPI types for Go backend: (Tool to be identified, e.g., `oapi-codegen`) - `internal/generated/types.go`
    - OpenAPI types for TypeScript frontend: (Tool to be identified, e.g., `openapi-typescript-codegen`) - `frontend/src/generated/api-types.ts`

## Technical Constraints

- (To be defined as they arise)

## Dependencies

- Key backend dependencies will be listed in `go.mod`.
- Key frontend dependencies will be listed in `frontend/package.json`.

## Tool Usage Patterns

- **API Development:** Define API in OpenAPI, generate types/server stubs, implement business logic.
- **Database Migrations:** Create `up` and `down` SQL files for schema changes.
- **Frontend Development:** Component-based architecture, TypeScript for type safety, Vite for fast development.

*(This is an initial assessment based on the file structure. It will be refined as more information is gathered by reading specific configuration files and code.)*
