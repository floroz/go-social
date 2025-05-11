# Active Context

## Current Work Focus

- **Planning Next Steps:** Discussing the next chunks of work with the user after successfully implementing and verifying the frontend signup payload wrapping.

## Recent Changes

- **Completed Memory Bank Update (Previous Task):** Updated all memory bank files to reflect the resolution of `TestUserSignup_ValidationErrors` functional test failures.
- **Implemented Frontend Signup Payload Wrapping (Chunk A.3):**
    - **Frontend Types (`frontend/src/generated/api-types.ts`):** Confirmed that generated types correctly expect a wrapped request payload (`{ data: SignupRequest }`) for the signup operation, following OpenAPI updates from Chunk A.1.
    - **Authentication Service (`frontend/src/services/authService.ts`):**
        - Modified the `signup` method to wrap the `signupData` (which is the flat `SignupRequest` type) within a `data` object before sending it in the `apiClient.post` call. The method signature remains `signup(signupData: SignupRequest)`.
        - Example change: `apiClient.post("/v1/auth/signup", { data: signupData });`
    - **Verification:**
        - User confirmed that the signup form on the frontend now works correctly, and the request payload sent to the backend is appropriately wrapped.

## Next Steps

1.  Discuss the next chunks of work with the user, likely focusing on:
    *   Further backend test coverage improvements (remainder of Chunk A.2).
    *   Rolling out the payload wrapping convention to other API endpoints (Part B).
    *   Or, addressing other items from "What's Left to Build" in `progress.md`.
2.  Update memory bank files based on the decided next steps.

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
