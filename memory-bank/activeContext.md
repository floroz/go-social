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

1. Update `memory-bank/systemPatterns.md` and `memory-bank/progress.md` to reflect the revised plan for enforcing wrapped request bodies.
2. Present the revised technical plan (which involves changing OpenAPI spec and frontend, keeping backend handler as-is for request structure) to the user.
3. Await user approval before implementing this new plan.

## Active Decisions and Considerations

- **Revised Payload Convention Strategy (User Directive):** The project aims to enforce a consistent `{"data": <payload>}` structure for **both request bodies and success response bodies**. This is a change from the previous understanding where only responses were wrapped.
- **Impact on Signup Error Fix:**
    - The previous plan was to change the backend handler to accept a flat *request* payload, aligning with the *current* OpenAPI spec.
    - The **new plan** is to:
        1.  Modify the OpenAPI specification for `/v1/auth/signup` to define a *wrapped* request body (e.g., `{"data": {"first_name": ...}}`).
        2.  Regenerate frontend types based on the new OpenAPI spec.
        3.  Update the frontend to send this wrapped request payload.
        4.  The backend handler (`cmd/api/auth_handlers.go`), which already expects this wrapped structure for requests, would then work correctly without modification to its unmarshaling logic.
- This new approach makes the backend handler's current expectation for a wrapped request correct, and the OpenAPI spec/frontend the parts that need to change to achieve consistency.
- The field name `first_name` (snake_case) remains consistent. The issue is purely structural (flat vs. wrapped request).

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
- **Revised understanding of signup error cause & fix:**
    - **Initial Error:** `json: unknown field "first_name"` when backend `signupHandler` (expecting `{"data": {...}}` for request) received a flat request `{"first_name": ...}` from frontend (which was aligned with the then-current OpenAPI spec for requests).
    - **New Directive:** Enforce `{"data": {...}}` wrapper for *all request bodies* for consistency with response wrappers.
    - **Revised Root Cause of Incompatibility:** The *current OpenAPI specification* for signup requests (and other requests) defines them as flat, which is now considered out of sync with the new desired project convention. The backend handler's expectation of a wrapped request is now the *target state*.
    - **Revised Fix Strategy:** Modify OpenAPI spec for signup request to be wrapped, regenerate frontend types, and update frontend to send wrapped payload. The backend handler's request unmarshaling logic will then be correct.

*(This file will be updated frequently as work progresses.)*
