# System Patterns

## System Architecture

The project appears to follow a standard client-server architecture:

- **Frontend (Client):** A Single Page Application (SPA) built with React and TypeScript, served by Vite. It interacts with the backend via HTTP API calls.
- **Backend (Server):** A Go application providing a RESTful API. It handles business logic and interacts with the database.
- **Database:** A PostgreSQL database stores persistent data.

```mermaid
graph TD
    User[User] -- Interacts via Browser --> Frontend[Frontend SPA (React/TS)]
    Frontend -- HTTP API Calls --> Backend[Backend API (Go)]
    Backend -- SQL Queries --> Database[Database (PostgreSQL)]
```

## Key Technical Decisions (Inferred)

- **API-First Approach:** The use of OpenAPI (`openapi.yaml`) suggests that the API contract is defined first, and then code (types, server stubs) is generated from this definition for both backend and frontend. This promotes consistency and clear separation of concerns.
- **Monorepo Structure (Potentially):** While not a strict monorepo in the sense of using tools like Lerna or Turborepo, the backend and frontend code reside in the same top-level repository. This can simplify coordinated changes but requires clear separation of concerns within the directory structure.
- **Layered Architecture (Backend):** The backend directory structure (`cmd/api`, `internal/services`, `internal/repositories`, `internal/domain`) suggests a layered architecture:
    - `cmd/api` (Handlers): Handles HTTP requests and responses, delegates to services.
    - `internal/services`: Contains business logic, orchestrates repositories.
    - `internal/repositories`: Abstracts data access, interacts with the database.
    - `internal/domain`: Defines core data structures (models) and business rules.
- **Component-Based UI (Frontend):** The use of React and the `frontend/src/components/` directory indicates a component-based approach to building the user interface.
- **Type Safety:** TypeScript on the frontend and Go's static typing on the backend, along with generated types from OpenAPI, aim to ensure type safety across the stack.

## Design Patterns in Use (Potential)

- **Repository Pattern (Backend):** `internal/repositories` strongly suggests its use for decoupling business logic from data access logic.
- **Service Layer Pattern (Backend):** `internal/services` indicates this pattern for encapsulating business logic.
- **Dependency Injection (Backend):** Go applications often use dependency injection (passing dependencies as arguments to functions or structs) to manage dependencies between layers (e.g., injecting a repository into a service). This will need to be confirmed by inspecting the code.
- **RESTful API Design:** Implied by the use of HTTP methods and resource-based URLs (to be confirmed by inspecting `openapi.yaml` and handler implementations).
- **State Management (Frontend):** `frontend/src/stores/authStore.ts` suggests a centralized or feature-specific state management pattern. The exact pattern (e.g., Flux, Redux-like, Zustand, Jotai) needs to be identified.
- **Mocking for Tests (Frontend):** `frontend/src/mocks/` and MSW usage indicate a pattern of mocking API responses for frontend testing.

## Component Relationships (High-Level)

- **Authentication:**
    - Frontend: `LoginPage.tsx`, `SignupPage.tsx`, `useLogin.ts`, `useSignup.ts`, `authService.ts`, `authStore.ts`.
    - Backend: `auth_handlers.go`, `auth_service.go`, likely interacting with `user_repository.go`.
- **Posts:**
    - Frontend: (Components to be created/identified)
    - Backend: `post_handlers.go`, `post_service.go`, `post_repository.go`, `post_model.go`.
- **Comments:**
    - Frontend: (Components to be created/identified)
    - Backend: `comments_handlers.go`, `comment_service.go`, `comment_repository.go`, `comment_model.go`.

## Critical Implementation Paths

- **User Authentication Flow:** Registration, login, session management (e.g., JWTs, cookies).
- **Post Creation and Retrieval:** How posts are created, stored, and displayed.
- **Comment Creation and Retrieval:** How comments are associated with posts and users.
- **API Request/Response Cycle:** How data flows from frontend to backend and back, including error handling.
    - **Payload Structure (Signup):** Investigation of a signup error (Nov 2025) revealed that the OpenAPI specification and frontend implementation correctly use a flat JSON payload for signup (e.g., `{"first_name": "..."}`). However, the backend handler (`cmd/api/auth_handlers.go`) at that time incorrectly expected a nested payload (e.g., `{"data": {"first_name": "..."}}`). This discrepancy was the root cause of the `json: unknown field "first_name"` error. The field naming convention itself (snake_case, e.g., `first_name`) was consistent. This highlights the importance of ensuring backend handlers strictly adhere to the OpenAPI contract.

*(This is an initial assessment based on the file structure and common practices. It will be refined by examining the code and configuration files in more detail.)*
