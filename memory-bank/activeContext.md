# Active Context

## Current Work Focus

- **Implementing Home Page (Phase 2 - Feed Display):** Building the core components for displaying the post feed on the Home Page, starting with `PostService`.

## Recent Changes

- **Corrected Phase 1 Complete (Auth & Routing Foundation for Home Page):**
    - Re-created `frontend/src/pages/HomePage.tsx` as a placeholder.
    - Re-created `frontend/src/components/ProtectedRoute.tsx` to handle auth checks and redirection.
    - Modified `frontend/src/main.tsx` to use `ProtectedRoute` for the `/` route, rendering `HomePage`.
    - Updated `frontend/src/pages/SignupPage.tsx` to include a "Log in" link.
- **Reverted Incorrect Home Page Phase 1 Implementation:**
    - Restored `frontend/src/App.tsx` to its original state (pre-routing changes).
    - Reverted `frontend/src/pages/SignupPage.tsx` by removing the "Login" link and `Link` import.
    - Deleted `frontend/src/pages/HomePage.tsx` and `frontend/src/components/ProtectedRoute.tsx` as they were part of the incorrect iteration.
    - **Mistake Identified:** Attempted to implement routing in `App.tsx` instead of the existing setup in `frontend/src/main.tsx`. Proceeded with implementation steps without explicit user go-ahead after plan presentation.
- **Updated `techContext.md`:** Corrected information about `react-router` v6.4+ (includes `react-router-dom` functionalities based on user feedback).
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
    - **Frontend Service Updates (Ongoing):**
        - `frontend/src/services/authService.ts`: Updated `login` method to send wrapped payload; confirmed by user. (Signup was done previously).
        - **User Profile Update (`PUT /v1/users`):** Searched frontend code (`frontend/src`) for usage of this endpoint. Found it only in generated types (`frontend/src/generated/api-types.ts`). This indicates the "update user profile" feature is not yet implemented in the frontend. No code changes required for this endpoint at this time.
        - **Post Creation/Update (`POST /v1/posts`, `PUT /v1/posts/{id}`):** Searched frontend code (`frontend/src`) for usage of these endpoints. Found them only in generated types (`frontend/src/generated/api-types.ts`). This indicates the "create/update post" features are not yet implemented in the frontend. No code changes required for these endpoints at this time.
        - **Comment Creation/Update (`POST /v1/posts/{postId}/comments`, `PUT /v1/posts/{postId}/comments/{id}`):** Searched frontend code (`frontend/src`) for usage of these endpoints. Found them only in generated types (`frontend/src/generated/api-types.ts`). This indicates the "create/update comment" features are not yet implemented in the frontend. No code changes required for these endpoints at this time.
    - **Conclusion of Frontend Convention Rollout (for existing features):** All relevant V1 POST/PUT endpoints have been checked. Only `authService.ts` (login/signup) required modifications, which were completed. Other features (user profile, posts, comments) appear unimplemented in the frontend.

## Next Steps

**Corrected Phase 1: Authentication and Routing Foundation (Home Page Plan) - COMPLETED**
1.  **Re-create `frontend/src/pages/HomePage.tsx`** (as a placeholder initially). - COMPLETED
2.  **Re-create `frontend/src/components/ProtectedRoute.tsx`**. - COMPLETED
3.  **Modify `frontend/src/main.tsx` for Protected Home Route**. - COMPLETED
4.  **Update `frontend/src/pages/SignupPage.tsx`:** Add "Already have an account? Login" link. - COMPLETED

**Phase 2: Home Page - Feed Display (Home Page Plan) - CURRENT**
1.  Create `PostService` (`frontend/src/services/postService.ts`) with `listPosts()` and `createPost()` methods. - **NEXT IMMEDIATE STEP**
2.  Create `usePosts` Hook (`frontend/src/hooks/usePosts.ts`) for fetching and managing post state.
3.  Create `PostCard` Component (`frontend/src/components/PostCard.tsx`) to display individual posts.
4.  Integrate post fetching and display into `HomePage.tsx`.
5.  Update `memory-bank/progress.md` as these steps are completed.

**Phase 3: Basic Comment Interaction (Home Page Plan - Future Enhancement Path) - PENDING**
*   (Tasks remain the same)

## Active Decisions and Considerations

**Home Page Implementation Plan (Nov 5, 2025):**
*   **Goal:** Create a Home Page displaying a feed of posts for logged-in users. Enhance signup/login flow.
*   **Authentication:** Home Page (`/`) will be a protected route, redirecting to `/login` if the user is not authenticated. `authStore.ts` will manage auth state.
*   **Signup Enhancement:** `SignupPage.tsx` will include a link to `/login`.
*   **Feed Components:**
    *   `HomePage.tsx`: Main page component.
    *   `postService.ts`: Service for `GET /v1/posts` (and other post actions).
    *   `usePosts.ts`: Custom hook for post data fetching and state.
    *   `PostCard.tsx`: Component to display individual post details (author, content, timestamp, comment/like counts placeholders, view/add comment links).
*   **API Usage:** Primarily `GET /v1/posts` (requires auth) for the feed. Comment interactions (`GET /v1/posts/{postId}/comments`, etc.) will also require auth.
*   **Initial Comment Handling:** Likely navigate to a separate `PostDetailPage` for full comment interaction to simplify the initial Home Page.
*   **Styling:** Utilize Shadcn/UI components for consistency.
*   **Future Considerations:** Pagination/infinite scroll for posts, inline "Create Post" form, real-time updates.

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
    - **Router Setup Location:** Always verify `main.tsx` (or the application's entry point) for existing router configurations before attempting to add or modify routing in components like `App.tsx`.

*(This file will be updated frequently as work progresses.)*
