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

1. Update `memory-bank/systemPatterns.md` and `memory-bank/progress.md` to reflect the latest revised plan (inline OpenAPI wrappers, testing, rollout strategy).
2. Present the latest revised technical plan to the user for feedback.
3. Await user approval before any implementation.

## Active Decisions and Considerations

- **Payload Convention Enforcement (User Directive):**
    - All API request bodies and success response bodies must use a `{"data": <payload>}` wrapper.
    - Error responses must use an `{"errors": [...]}` wrapper.
    - **OpenAPI Definition Style:** The `data` wrapper for request bodies should be defined *inline* within the OpenAPI path definition's `requestBody.content.application/json.schema`, rather than creating separate named wrapper schemas (e.g., `WrappedSignupRequest`).
- **Revised Plan for Signup Error & Broader Rollout:**
    - **Part A (Signup Endpoint):**
        1.  Modify `openapi/v1/paths/auth.yaml` for `/v1/auth/signup` to use an inline `data` wrapper in the `requestBody` schema, referencing the existing flat `SignupRequest` schema for its `data` property.
        2.  Confirm `SignupSuccessResponse` already uses the `data` wrapper (it does).
        3.  Regenerate frontend types.
        4.  Update frontend `SignupPage.tsx` to send the wrapped request.
        5.  Verify backend `signupHandler` (which expects wrapped request) now works correctly.
        6.  Review and add backend tests for signup request unmarshaling (wrapped) and response structure (wrapped).
    - **Part B (Rollout to Other Endpoints):**
        1.  Identify other POST/PUT/PATCH endpoints with flat request bodies.
        2.  Iteratively apply the same process: update OpenAPI with inline `data` wrapper for requests, regenerate types, update frontend, update backend handler (if it currently expects flat requests), and add/update tests for both request and response structures.
- This strategy ensures consistency and makes the backend's current expectation for a wrapped signup request the correct target, requiring changes primarily in the OpenAPI spec and frontend.

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
- It utilizes OpenAPI for API design and code generation.
- Docker is used for containerization.
- A comprehensive set of `.clinerules` dictates coding standards and best practices for TypeScript development.
- **Latest understanding of signup error fix & project conventions:**
    - **User Directive:** Enforce `{"data": {...}}` wrapper for all request bodies and success responses. Error responses: `{"errors": [...]}`. Request wrappers in OpenAPI are to be defined inline in path definitions.
    - **Path to Fix Signup Error:**
        1.  Modify `openapi/v1/paths/auth.yaml` to define an inline `data` wrapper for the `/v1/auth/signup` request body, which will contain the properties of the original `SignupRequest`.
        2.  Regenerate frontend types.
        3.  Update `SignupPage.tsx` to send this wrapped payload.
        4.  The backend `signupHandler` already expects this wrapped structure for requests, so its unmarshaling logic should then work correctly.
        5.  Add/verify backend tests for request and response shapes.
    - **Broader Impact:** This convention will be rolled out to other relevant endpoints.

*(This file will be updated frequently as work progresses.)*
