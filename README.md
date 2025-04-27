# Go-Social

[![CI Pipeline](https://github.com/floroz/go-social/actions/workflows/ci.yml/badge.svg)](https://github.com/floroz/go-social/actions/workflows/ci.yml)

Go-Social is a social media application built with Go (backend) and React (frontend). It provides features for users to create posts, comment on posts, and interact with each other. The backend uses PostgreSQL as the database and follows a clean architecture.

## Tech Stack

*   **Backend:** Go, Chi Router, PostgreSQL
*   **Frontend:** React, TypeScript, Vite, Tailwind CSS, Shadcn/ui, Zustand, React Query, Vitest, MSW
*   **Database Migrations:** golang-migrate
*   **Development:** Docker, Air (optional, for Go live reload)

## Getting Started

### Prerequisites

*   Go (version specified in `go.mod`)
*   Node.js (version specified in `frontend/package.json` or latest LTS)
*   npm (comes with Node.js)
*   Docker & Docker Compose
*   [Golang Migrate CLI](https://github.com/golang-migrate/migrate?tab=readme-ov-file#cli-usage)
*   [Air](https://github.com/air-verse/air) (Optional, for Go backend live reload)

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/floroz/go-social.git
    cd go-social
    ```

2.  **Backend Dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Frontend Dependencies:**
    ```sh
    cd frontend
    npm install
    cd ..
    ```

4.  **Environment Variables:**
    *   Copy `.env.local.example` (if it exists) to `.env.local` and configure backend variables (DB connection, JWT secret).
    *   Copy `frontend/.env.example` (if it exists) to `frontend/.env.development` and `frontend/.env.production` and configure frontend variables (mainly `VITE_API_BASE_URL`). Ensure the development URL matches the backend setup (e.g., `http://localhost:8080/api`).

5.  **Start Database:**
    ```sh
    docker compose up -d
    ```

6.  **Run Database Migrations:**
    ```sh
    make migrate-up
    ```

## Development

### Running the Application

To start both the backend and frontend development servers concurrently:

```sh
make dev
```

This will:
*   Start the Go backend (using `go run`) in the background on port 8080 (or as configured).
*   Start the Vite frontend dev server in the foreground on port 5173 (or the next available port).

Alternatively, you can run them separately:

*   **Backend Only (with live reload using Air):**
    ```sh
    air
    ```
*   **Backend Only (standard Go run):**
    ```sh
    make dev-be
    # or
    go run ./cmd/main.go
    ```
*   **Frontend Only:**
    ```sh
    make dev-fe
    # or
    cd frontend && npm run dev
    ```

## Testing

*   **Run Backend Tests:**
    ```sh
    make test
    ```
*   **Run Frontend Tests:**
    ```sh
    make test-fe
    # or
    cd frontend && npm run test
    ```
*   **Run All Tests (Backend & Frontend):**
    ```sh
    make test-all
    ```

## API Spec

This project uses an OpenAPI 3 specification as the single source of truth for the API contract. This drives code generation for both backend and frontend types, ensuring consistency.

### Response Structure Convention

To maintain consistency across all endpoints, API responses adhere to the following structure:

*   **Success Responses:** Contain a top-level `data` key holding the actual response payload.
    ```json
    // Example: 201 Created from /v1/auth/signup
    {
      "data": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        // ... other user fields
      }
    }
    ```
    ```json
    // Example: 200 OK from /v1/auth/login
    {
      "data": {
        "token": "eyJ..."
      }
    }
    ```
*   **Error Responses:** Contain a top-level `errors` key holding an array of error objects. Each error object includes a `code` and `message`, and optionally a `field` for validation errors.
    ```json
    // Example: 400 Bad Request
    {
      "errors": [
        {
          "code": "VALIDATION_ERROR",
          "message": "Email format is invalid.",
          "field": "email"
        },
        {
          "code": "VALIDATION_ERROR",
          "message": "Password must be at least 8 characters.",
          "field": "password"
        }
      ]
    }
    ```
    ```json
    // Example: 401 Unauthorized
    {
      "errors": [
        {
          "code": "UNAUTHORIZED",
          "message": "Invalid email or password."
        }
      ]
    }
    ```

### Specification Files Structure

The OpenAPI specification files are located in the `openapi/` directory:

*   `openapi/openapi.yaml`: The main entry point, defining info, servers, tags, and references to paths and shared schemas.
*   `openapi/shared/`: Contains schemas (`common.yaml`, `user.yaml`, etc.) shared across API versions or endpoints. This includes the standard `ApiError` and `ApiErrorResponse` schemas.
*   `openapi/v1/`: Contains definitions specific to V1 of the API, including paths (`paths/auth.yaml`, etc.) and version-specific schemas (`schemas/auth.yaml`, etc.). This includes success response wrappers like `SignupSuccessResponse` and `LoginSuccessResponse`.

### Code Generation

We use code generation tools to create Go types for the backend and TypeScript types for the frontend based on the OpenAPI specification.

*   **Tools:**
    *   `oapi-codegen`: Generates Go types (`internal/generated/types.go`).
    *   `openapi-typescript`: Generates TypeScript types (`frontend/src/generated/api-types.ts`).
    *   `@redocly/cli`: Used to bundle the multi-file OpenAPI spec into a single file (`openapi/openapi-bundled.yaml`) before generation.
*   **Tool Management:**
    *   Go tools (`oapi-codegen`, `golangci-lint`) are managed via `go.mod` using a `tools.go` file. Run `go install` for the packages listed in `tools.go` if needed (or add a `make setup-tools` target).
    *   Node.js tools (`@redocly/cli`, `openapi-typescript`) are managed as dev dependencies in `frontend/package.json`. Run `npm install` within the `frontend` directory.
*   **Generating Types:**
    *   To regenerate both Go and TypeScript types after modifying the OpenAPI spec, run:
        ```sh
        make generate-types
        ```
    *   This command first bundles the spec using `redocly` and then runs `oapi-codegen` and `openapi-typescript`.

### Facade Pattern

To decouple the main application code from the potentially verbose or unstable generated code, we use facade modules:

*   **Backend:** `internal/apitypes/types.go` re-exports the necessary types from `internal/generated`. Backend handlers should import from `internal/apitypes`, not directly from `internal/generated`.
*   **Frontend:** `frontend/src/types/api.ts` re-exports the necessary types from `frontend/src/generated/api-types.ts`. Frontend code (services, components) should import from `frontend/src/types/api`, not directly from `frontend/src/generated`.

### Accessing API Documentation (Swagger UI)

The backend serves an interactive Swagger UI for exploring the API based on the OpenAPI specification. Once the backend server is running (e.g., via `make dev` or `make dev-be`), you can access the documentation at:

[http://localhost:8080/api/docs](http://localhost:8080/api/docs)

(Adjust the port if your backend runs on a different one).
