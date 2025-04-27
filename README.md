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
