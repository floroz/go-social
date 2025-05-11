# Active Context

## Current Work Focus

- **Finalizing Memory Bank Update:** Reflecting completion of the backend/functional test phase for the request payload convention rollout across all relevant V1 POST/PUT endpoints. Preparing to discuss the next phase of work.

## Recent Changes

- **Completed Memory Bank Update (Previous Task):** Updated memory bank after fixing `TestUserSignup_ValidationErrors`.
- **Implemented Frontend Signup Payload Wrapping (Chunk A.3 - Completed):**
    - Ensured frontend sends wrapped `{"data": ...}` for signup; confirmed by user.
- **Completed Backend/Functional Test Rollout of Request Payload Convention (Part B):**
    - Systematically reviewed all V1 API endpoints (`auth.yaml`, `user.yaml`, `post.yaml`, `comment.yaml`).
    - For all POST/PUT endpoints requiring a request body, the following steps were completed:
        1.  **OpenAPI Specification:** Updated to define request bodies with an inline `{"data": <ActualPayloadSchema>}` wrapper.
        2.  **Type Generation:** Ran `make generate-types` successfully after each OpenAPI modification.
        3.  **Backend Handlers:** Updated or verified handlers in `cmd/api/` to correctly unmarshal the wrapped request payloads.
        4.  **Functional Tests:** Updated relevant tests in `test/functional/` to send wrapped request payloads.
        5.  **Verification:** Ensured `make test` passed after each set of backend/test modifications.
    - **Endpoints covered:**
        - `/v1/auth/login` (request body)
        - `PUT /v1/users` (request body)
        - `POST /v1/posts` (request body)
        - `PUT /v1/posts/{id}` (request body)
        - `POST /v1/posts/{postId}/comments` (request body)
        - `PUT /v1/posts/{postId}/comments/{id}` (request body)
    - **Frontend Service Updates (Partial):**
        - `frontend/src/services/authService.ts`: Updated `login` method to send wrapped payload; confirmed by user. (Signup was done previously).
        - Frontend updates for user, post, and comment services are pending.

## Next Steps

1.  Update `memory-bank/progress.md` to reflect the completion of the backend/functional test convention rollout.
2.  Discuss the next phase with the user:
    *   **Backend Test Enhancement (Signup - Remainder of Chunk A.2):** Detailed review and enhancement of backend tests for the signup endpoint for comprehensive coverage.
    *   **Frontend Updates for Convention Rollout:** Update frontend services for user, post, and comment operations to send wrapped request payloads.
    *   Address other items from "What's Left to Build" in `progress.md`.

## Active Decisions and Considerations

- **Critical Quality Steps Added to Workflow:**
    - After any OpenAPI specification change: **Must run `make generate-types`**.
    - After any backend code change (including test modifications or handler adjustments, or post-type-generation verification): **Must run `make test`**.
- **Payload Convention Enforcement (User Directive - Confirmed):**
    - All API request bodies and success response bodies: `{"data": <payload>}` wrapper.
    - Error responses: `{"errors": [{"code": "...", "field": "...", "message": "..."}]}` wrapper.
    - OpenAPI request wrappers: Defined *inline* in path definitions.
- **API Error Handling Strategy:**
    - Service layer returns detailed validation errors (e.g., `validator.ValidationErrors`).
    - API layer (`handleErrors`) processes these, extracts relevant information (like field names, converting to appropriate case for API response), and formats a standardized error response. This keeps service layer errors more Go-idiomatic and centralizes API error formatting.
- **Iterative Development:** The approach of fixing tests by modifying service, then API handlers, then re-testing, proved effective.

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
    - **API Error Handling (Refined Nov 2025):**
        - Services should return specific error types or `validator.ValidationErrors` to convey detailed error information.
        - The central `handleErrors` function in `cmd/api/api.go` is responsible for:
            - Type-asserting errors to known types (e.g., `domain.ErrDuplicateEmailOrUsername`, `validator.ValidationErrors`, custom `*domain.Error` types).
            - Extracting necessary details (e.g., field names from `validator.FieldError`, specific messages).
            - Converting field names to the API's desired case (e.g., snake_case using `toSnakeCase` helper).
            - Constructing the standardized `apitypes.ApiErrorResponse` with appropriate HTTP status codes, API error codes, messages, and field names.
        - This pattern centralizes the translation of internal Go errors into user-facing API error responses.
    - **Payload Conventions:** The `{"data": ...}` wrapper for requests/success responses and `{"errors": ...}` for error responses remains a key convention.

*(This file will be updated frequently as work progresses.)*
