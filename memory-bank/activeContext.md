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

1. Update other relevant Memory Bank files (`systemPatterns.md`, `progress.md`) with findings.
2. Present the technical plan to fix the signup error to the user.
3. Await user approval before implementing the fix.

## Active Decisions and Considerations

- The primary cause of the "json: unknown field 'first_name'" error is the backend `signupHandler` expecting a `{"data": ...}` wrapper in the JSON payload, while the OpenAPI specification and the frontend implementation correctly use a flat payload structure. The field name `first_name` itself is consistent (snake_case) between the frontend payload, OpenAPI spec, and the backend's `CreateUserDTO` struct tags.

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
- It utilizes OpenAPI for API design and code generation. The frontend types are generated from this spec.
- Docker is used for containerization.
- A comprehensive set of `.clinerules` dictates coding standards and best practices for TypeScript development.
- **Key finding on signup error:** The backend `signupHandler` in `cmd/api/auth_handlers.go` incorrectly expects a nested `{"data": ...}` payload for signup, contradicting the OpenAPI specification and the frontend's (correct) implementation of a flat payload. The field name `first_name` (snake_case) is consistent across the OpenAPI spec, generated frontend types, and the backend `CreateUserDTO`'s JSON tags. The error arises because the flat payload's `first_name` field is "unknown" to the top-level struct the backend handler initially tries to unmarshal into (which expects a `data` field).

*(This file will be updated frequently as work progresses.)*
