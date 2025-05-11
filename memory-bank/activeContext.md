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

1. Update `memory-bank/systemPatterns.md` to further clarify request vs. response payload conventions.
2. Re-present the confirmed technical plan to fix the signup error to the user, incorporating the clarification about payload conventions.
3. Await user approval before implementing the fix.

## Active Decisions and Considerations

- **Confirmed Payload Conventions:** Further investigation (prompted by user feedback) into `openapi/v1/schemas/*.yaml` files (auth, post, user) confirms:
    - **Request Bodies:** OpenAPI specifications for signup, post creation, and user updates consistently define *flat* request body schemas (e.g., `SignupRequest`, `CreatePostRequest` are flat objects).
    - **Success Response Bodies:** OpenAPI specifications consistently define success responses with a `{"data": ...}` wrapper.
- **Validity of Original Plan:** The original plan to modify the backend `signupHandler` to accept a flat request payload remains correct. This aligns the handler with the OpenAPI contract for *request bodies* and fixes the "json: unknown field 'first_name'" error. The backend handler's current expectation of a `{"data": ...}` wrapper for the *request* is the actual deviation.
- The field name `first_name` itself is consistent (snake_case) between the frontend payload, OpenAPI spec, and the backend's `CreateUserDTO` struct tags.

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
- **Key finding on signup error (re-confirmed):** The backend `signupHandler` in `cmd/api/auth_handlers.go` incorrectly expects a nested `{"data": ...}` payload for signup. This contradicts the OpenAPI specification, which defines a *flat* payload for signup requests (and other POST/PUT requests like post creation and user update). The project's convention of using a `{"data": ...}` wrapper applies to *success responses*, not request bodies, as per the OpenAPI definitions. The field name `first_name` (snake_case) is consistent. The error arises because the flat payload's `first_name` field is "unknown" to the top-level struct the backend handler initially tries to unmarshal the request into.

*(This file will be updated frequently as work progresses.)*
