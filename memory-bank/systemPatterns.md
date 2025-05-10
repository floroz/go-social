# System Patterns

## System Architecture

- **Monolithic Backend:** Go application serving API endpoints.
- **Single Page Application (SPA) Frontend:** React application consuming the API.
- **Database:** PostgreSQL for data persistence.
- **API Layer:** Defined by OpenAPI, acts as the contract between backend and frontend.

## Key Technical Decisions

- **OpenAPI-first approach:** The API specification (`openapi.yaml` and partials) is the source of truth.
- **Code Generation:** Go types (`internal/generated/types.go`) and TypeScript types (`frontend/src/generated/api-types.ts`) are generated from the OpenAPI spec. This ensures consistency and reduces manual type definition.
- **API Versioning:** Paths are prefixed with `/v1/`, allowing for future API evolution without breaking existing clients.
- **Structured Error Responses:** Consistent error format (`{"errors": [{"code": "...", "message": "...", "field": "..."}]}`) for all API endpoints.
- **Structured Success Responses:** Consistent success format (`{"data": {...}}`) for all API endpoints.
- **Use of Partials for OpenAPI:** Specification is broken down into smaller, manageable YAML files (in `openapi/shared/` and `openapi/v1/`) and bundled into `openapi/openapi-bundled.yaml` using `@redocly/cli`.

## Design Patterns in Use

- **Facade Pattern:**
    - Backend: `internal/apitypes/types.go` re-exports types from `internal/generated`, decoupling service logic from generated code.
    - Frontend: `frontend/src/types/api.ts` re-exports types from `frontend/src/generated/api-types.ts`, decoupling components/services from generated code.
- **Repository Pattern (implied):** The backend likely uses repositories for database interactions (e.g., `internal/repositories/`).
- **Service Layer (implied):** Business logic is likely encapsulated in services (e.g., `internal/services/`).
- **Middleware (explicit):** Chi router is used, which supports middleware for concerns like authentication (`cmd/middlewares/auth_middleware.go`).

## Component Relationships

- **Frontend (`frontend/`)** depends on **Backend API**.
- **Backend API (`cmd/api/`)** depends on **Services (`internal/services/`)**.
- **Services (`internal/services/`)** depend on **Repositories (`internal/repositories/`)** and **Domain Models (`internal/domain/`)**.
- **Repositories (`internal/repositories/`)** depend on **Database (`cmd/database/db.go`)** and **Domain Models (`internal/domain/`)**.
- **OpenAPI Specification (`openapi/`)** is used by code generation tools to produce types for both **Backend (`internal/generated/`)** and **Frontend (`frontend/src/generated/`)**.

## Critical Implementation Paths

- **Authentication Flow:** Signup, login, logout, token refresh. Involves JWTs and secure handling of credentials.
- **Post Creation/Management:** Creating, reading, updating, deleting posts.
- **Comment Creation/Management:** Creating, reading, updating, deleting comments on posts.
- **API Spec Modification and Type Regeneration:** Any change to an API requires updating the OpenAPI YAML files, re-bundling, and re-generating Go and TypeScript types using `make generate-types`.
