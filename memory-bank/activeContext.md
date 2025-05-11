# Active Context

## Current Work Focus

- **Investigate and Plan Fix for Signup 400 Error:** Diagnosing a 400 error (`json: unknown field "first_name"`) on the `/v1/auth/signup` endpoint. The current phase involves investigation and proposing a technical plan.

## Recent Changes

- **Completed Memory Bank Refresh:** Read all core memory bank files.
- **Investigated Signup Error - Frontend:**
    - Read `frontend/src/pages/SignupPage.tsx`: Confirmed form uses `firstName` (camelCase) internally but maps to `first_name` (snake_case) for the `SignupRequest` payload.
    - Read `frontend/src/types/api.ts`: Confirmed `SignupRequest` is an alias for the generated `components["schemas"]["SignupRequest"]`.
    - Read `frontend/src/generated/api-types.ts`: Confirmed generated `SignupRequest` schema uses `first_name` and `last_name` (snake_case), indicating OpenAPI spec uses snake_case for these fields in the payload.
- **Investigated Signup Error - Backend & OpenAPI:**
    - Read `cmd/api/auth_handlers.go`: Found that `signupHandler` expects the request body to be nested under a `data` key (e.g., `{"data": {"first_name": ...}}`), unmarshaling into `*domain.CreateUserDTO`.
    - Read `internal/domain/user_model.go`: Confirmed `domain.CreateUserDTO` (via embedded `EditableUserField`) correctly uses `json:"first_name"` and `json:"last_name"` (snake_case).
    - Read `openapi/openapi.yaml` and `openapi/v1/paths/auth.yaml`: Confirmed the OpenAPI specification for `/v1/auth/signup` defines the request body schema directly as `SignupRequest` (i.e., a flat payload, not nested under `data`).

## Next Steps

1. Update `memory-bank/progress.md` to reflect completion of the revised Chunk A.1.
2. Present completion of revised Chunk A.1 to the user and await feedback/approval to proceed to Chunk A.2 (Backend Test Analysis & Adjustments for Signup Endpoint).

## Active Decisions and Considerations

- **Critical Quality Steps Added to Workflow:**
    - After any OpenAPI specification change: **Must run `make generate-types`**.
    - After any backend code change (including test modifications or handler adjustments, or post-type-generation verification): **Must run `make test`**.
- **Payload Convention Enforcement (User Directive - Confirmed):**
    - All API request bodies and success response bodies: `{"data": <payload>}` wrapper.
    - Error responses: `{"errors": [...]}` wrapper.
    - OpenAPI request wrappers: Defined *inline* in path definitions.
- **Revised Chunked, Backend-First Iterative Plan for Signup Endpoint (Part A):**
    - **Chunk A.1: Update OpenAPI, Regenerate Types, Verify Handler & Run Initial Backend Tests (COMPLETED)**
        1.  Modified `openapi/v1/paths/auth.yaml` for `/v1/auth/signup` `requestBody` to use an inline `data` wrapper referencing `SignupRequest`. (DONE)
        2.  Ran `make generate-types` to propagate OpenAPI changes. (DONE)
        3.  Verified existing backend `signupHandler` in `cmd/api/auth_handlers.go` still correctly expects the wrapped request structure for unmarshaling (post type-generation). (DONE)
        4.  Confirmed `SignupSuccessResponse` in OpenAPI already uses the `data` wrapper. (DONE)
        5.  Ran `make test`; all existing tests passed. (DONE)
    - **Chunk A.2: Backend Test Analysis & Adjustments for Signup Endpoint (PENDING USER APPROVAL)**
        1.  Analyze results from `make test` in A.1 (tests passed, but need deeper look at coverage for response shapes).
        2.  Update/add tests for signup to ensure they send wrapped requests (if applicable and not already covered) and, critically, validate that success responses are `{"data": <User>}` and error responses are `{"errors": [...]}`.
        3.  Adjust signup handler's response generation if it doesn't already produce correctly wrapped responses. Iterate with `make test` until tests pass.
    - **(User Feedback Point after Chunk A.2)**
    - **Chunk A.3: Frontend Implementation for Signup (Details later)**
        1.  Regenerate frontend types (already done in A.1.2, but re-confirm or re-run if any backend types changed that affect frontend).
        2.  Update `SignupPage.tsx` to send wrapped request.
        3.  Test frontend.
    - **(User Feedback Point after Chunk A.3)**
- **Rollout to Other Endpoints (Part B - Future):** Will follow a similar chunked, backend-first, test-driven approach including the mandatory `make generate-types` and `make test` steps.
- This iterative strategy with integrated quality checks allows for focused backend stabilization.

## Important Patterns and Preferences (from `.clinerules/`)

*   **Generic Functions & `any`**: Use `any` inside generic function bodies if TypeScript cannot match runtime logic to type logic (e.g., conditional return types). Avoid `as <type>` casting generally.
*   **Default Exports**: Avoid default exports unless required by a framework (e.g., Next.js pages). Prefer named exports.
*   **Discriminated Unions**: Proactively use for modeling data with varying shapes (e.g., event types, fetching states) to prevent impossible states. Use `switch` statements for handling.
*   **Enums**: Do not introduce new enums. Retain existing ones. Use `as const` objects for enum-like behavior. Be mindful of numeric enum reverse mapping.
*   **`import type`**: Use `import type` for all type imports, preferably at the top level.
*   **Installing Libraries**: Use package manager commands (`pnpm add`, `yarn add`, `npm install`) to install the latest versions, rather than manually editing `package.json`.
*   **`interface extends`**: Prefer `interface extends` over `&` for modeling inheritance due to performance.
*   **JSDoc Comments**: Use for functions and types if behavior isn't self-evident. Use `@link` for internal references.
*   **Naming Conventions**:
    *   kebab-case for files (`my-component.ts`)
    *   camelCase for variables/functions (`myVariable`)
    *   UpperCamelCase (PascalCase) for classes/types/interfaces (`MyClass`)
    *   ALL_CAPS for constants/enum values (`MAX_COUNT`)
    *   `T` prefix for generic type parameters (`TKey`)
*   **`noUncheckedIndexedAccess`**: Be aware that if enabled, object/array indexing returns `T | undefined`.
*   **Optional Properties**: Use sparingly. Prefer `property: T | undefined` over `property?: T` if the property's presence is critical but its value can be absent.
*   **`readonly` Properties**: Use by default for object types to prevent accidental mutation. Omit only if genuinely mutable.
*   **Return Types**: Declare return types for top-level module functions (except JSX components).
*   **Throwing Errors**: Consider result types (`Result<T, E>`) for operations that might fail (e.g., parsing JSON) instead of throwing, unless throwing is idiomatic for the framework (e.g., backend request handlers).

## Learnings and Project Insights

- The project "Go Social" is a full-stack application with a Go backend and a React/TypeScript frontend.
- It utilizes OpenAPI for API design and code generation. Key commands: `make generate-types`, `make test`.
- Docker is used for containerization.
- A comprehensive set of `.clinerules` dictates coding standards and best practices for TypeScript development.
- **Latest Plan Summary (Signup Error Fix & Conventions with Quality Steps):**
    - **Convention:** `{"data": ...}` for requests & success responses; `{"errors": ...}` for errors. OpenAPI request wrappers: inline.
    - **Signup Fix (Backend First, Iterative with Testing):**
        1.  Update OpenAPI for signup request (inline `data` wrapper).
        2.  Run `make generate-types`.
        3.  Verify backend handler's request unmarshaling.
        4.  Run `make test` (initial run).
        5.  Analyze test results; update/add tests for wrapped request handling and structured `data`/`errors` responses. Adjust handler response generation if needed. Iterate with `make test`.
        6.  (Later) Update frontend.
    - **General Rollout:** Apply this iterative, test-focused, backend-first approach (including `make generate-types` and `make test` at appropriate points) to other endpoints.

*(This file will be updated frequently as work progresses.)*
