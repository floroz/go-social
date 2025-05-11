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

- **API-First Approach:** The use of OpenAPI (`openapi.yaml`) dictates the API contract. Code/types are generated from this (e.g., via `make generate-types` for frontend and potentially backend). This promotes consistency.
- **Development Workflow Integration:**
    - Changes to OpenAPI spec **must** be followed by running `make generate-types`.
    - Backend code changes (including tests or handler logic) **must** be followed by running `make test`.
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
    - **Payload Structure Conventions (Further Revised - Nov 2025):**
        - **User Directive:** Enforce consistent wrapper structures:
            - **Request Bodies:** Must be wrapped with a `data` key (e.g., `{"data": {"actual_payload..."}}`).
            - **Success Response Bodies:** Must be wrapped with a `data` key (e.g., `{"data": {"user_details..."}}`).
            - **Error Response Bodies:** Must be structured with an `errors` key (e.g., `{"errors": [{"code": "...", "message": "..."}]}`).
        - **OpenAPI Definition Style for Request Wrappers:** The `data` wrapper for request bodies is to be defined *inline* within the `requestBody.content.application/json.schema` of the OpenAPI path definition. This inline schema will have a `data` property that then `$ref`s the actual flat payload schema (e.g., `SignupRequest`). This avoids creating separate named wrapper schemas (like `WrappedSignupRequest`).
    - **Application to Signup Endpoint (Nov 2025):**
        - **Initial Discrepancy:** OpenAPI defined a flat signup request, frontend sent flat, but backend handler expected a wrapped (`{"data":...}`) request, causing an error.
        - **Resolution Strategy:**
            1. Modify `openapi/v1/paths/auth.yaml` to define an *inline* `data` wrapper for the signup request body, which internally references the flat `SignupRequest` schema.
            2. Regenerate frontend types.
            3. Update frontend to send the wrapped request.
            4. The backend handler's existing expectation for a wrapped request will then align with the new contract.
            5. Add/update backend tests for request and response shapes.
    - **General Rollout:** This approach (inline `data` wrapper in OpenAPI for requests, frontend updates, backend handler alignment, comprehensive testing) will be applied iteratively to other POST/PUT/PATCH endpoints to ensure project-wide consistency.
    - **Standardized API Error Responses (Updated Nov 2025):**
        - Error responses consistently use the `{"errors": [{"code": "...", "field": "...", "message": "..."}]}` structure.
        - Specific API error codes (e.g., `GOSOCIAL-005-CONFLICT`, `GOSOCIAL-006-VALIDATION_ERROR`) are used to differentiate error types.
        - For validation errors (`GOSOCIAL-006-VALIDATION_ERROR`), the `field` attribute in the error object is populated with the name of the problematic input field (e.g., "email", "password", "first_name"). This is achieved by:
            1. Services (e.g., `userService.Create`) returning `validator.ValidationErrors` directly when input DTO validation fails.
            2. The central API error handler (`cmd/api/api.go#handleErrors`) detecting `validator.ValidationErrors`.
            3. The handler extracting the struct field name (e.g., "FirstName") from the `validator.FieldError`, converting it to snake_case (e.g., "first_name") using a helper function (`toSnakeCase`), and using this as the `field` in the `apitypes.ApiError`.
            4. The error message from the validator is used as the `message` in `apitypes.ApiError`.
        - Conflict errors (e.g., duplicate email/username during signup) result in an HTTP 409 status and the `GOSOCIAL-005-CONFLICT` API error code.
        - General bad request errors due to malformed JSON in `auth_handlers.go` (e.g., during signup or login) result in an HTTP 400 status and `GOSOCIAL-001-BAD_REQUEST` or `GOSOCIAL-006-VALIDATION_ERROR` (note: `signupHandler` uses `CodeValidationError` for `readJSON` failures, `loginHandler` uses `CodeBadRequest`).

*(This is an initial assessment based on the file structure and common practices. It will be refined by examining the code and configuration files in more detail.)*
