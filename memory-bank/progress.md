# Project Progress

## What Works

- **Initial Project Structure:** The basic directory structure for a Go backend and a React/TypeScript frontend is in place.
- **Memory Bank Initialization:** The core Memory Bank files have been created.
- **API Error Handling (Signup - Validation & Conflicts):**
    - User signup validation errors (e.g., password too short, invalid email format) now correctly return HTTP 400 with API error code `GOSOCIAL-006-VALIDATION_ERROR` and include the specific problematic field name (e.g., "password", "email", "first_name") in the response.
    - User signup attempts with duplicate email/username now correctly return HTTP 409 with API error code `GOSOCIAL-005-CONFLICT`.
    - Functional test `TestUserSignup_ValidationErrors` and related tests for duplicate signups are passing.
- **Frontend Signup Payload:** The frontend (`authService.ts`) now correctly sends the signup request payload wrapped in a `{"data": ...}` object, aligning with the API specification and backend expectations. User has confirmed this functionality.
- **Frontend Login Payload:** The frontend (`authService.ts`) now correctly sends the login request payload wrapped in a `{"data": ...}` object for the `/v1/auth/login` endpoint, aligning with API specification and backend. User has confirmed this functionality.
- **User Profile Update (Backend):** `PUT /v1/users` endpoint expects and handles wrapped `{"data": ...}` request payloads. Functional tests updated and passing.
- **Post Endpoints (Backend):** `POST /v1/posts` and `PUT /v1/posts/{id}` endpoints expect and handle wrapped `{"data": ...}` request payloads. Functional tests updated and passing.
- **Comment Endpoints (Backend):** `POST /v1/posts/{postId}/comments` and `PUT /v1/posts/{postId}/comments/{id}` endpoints expect and handle wrapped `{"data": ...}` request payloads. Functional tests updated and passing.
- **User Profile Update (Frontend - `PUT /v1/users`):** Investigation revealed this feature is not yet implemented in the frontend. No code changes for payload wrapping were required for this endpoint.
- **Post Creation/Update (Frontend - `POST /v1/posts`, `PUT /v1/posts/{id}`):** Investigation revealed these features are not yet implemented in the frontend. No code changes for payload wrapping were required for these endpoints.
- **Comment Creation/Update (Frontend - `POST /v1/posts/{postId}/comments`, `PUT /v1/posts/{postId}/comments/{id}`):** Investigation revealed these features are not yet implemented in the frontend. No code changes for payload wrapping were required for these endpoints.

## What's Left to Build (High-Level)

This list is based on the initial `projectbrief.md` and common features for a social media platform. It will be refined as development progresses.

### Core Features
- **User Authentication:**
    - Backend: Implement registration, login (e.g., password hashing, session/token generation), logout, and potentially password recovery.
    - Frontend: Implement UI for registration, login, logout. Manage auth state.
- **Post Management:**
    - Backend: CRUD APIs for posts (Create, Read, Update, Delete).
    - Frontend: UI for creating, displaying, editing, and deleting posts.
- **Comment Management:**
    - Backend: CRUD APIs for comments, associated with posts and users.
    - Frontend: UI for adding and displaying comments on posts.
- **User Profiles:**
    - Backend: API to fetch user profile information (e.g., username, posts).
    - Frontend: UI to display user profiles.
- **Feed/Timeline:**
    - Backend: API to fetch a feed of posts (e.g., chronological, personalized).
    - Frontend: UI to display the post feed.

### Supporting Infrastructure
- **Database Schema:** Finalize and implement the database schema beyond the initial migrations (users, posts, comments).
- **API Endpoints:** Implement all API endpoints defined in `openapi.yaml`.
- **Frontend Routing:** Set up client-side routing for different pages/views.
- **Styling and UI Polish:** Develop a consistent and appealing visual design.
- **Error Handling:** Robust error handling on both frontend and backend.
- **Validation:** Input validation on both client and server sides.
- **Testing:**
    - Backend: Unit and integration tests for services and repositories.
    - Frontend: Unit, component, and E2E tests.
- **Deployment:** Setup CI/CD pipelines and deployment strategy.

## Current Status

- **Phase:** API Implementation & Testing / Planning.
- **Current Focus:** Frontend payload convention rollout for V1 POST/PUT endpoints is complete based on existing features. Awaiting user direction for the next task.
- **Blockers:** None.

## Known Issues

- **Signup Endpoint Payload Wrapping:**
    - **Status:** RESOLVED. Backend expects wrapped (`{"data": ...}`) requests, OpenAPI defines this, and frontend now sends the wrapped request.
- **Comprehensive Backend Testing for Response Shapes:** While the signup validation error response shape is now correct, a broader review and addition of tests for other endpoints' request/response shapes (ensuring `data` and `errors` wrappers are consistently tested) is an ongoing effort (related to remainder of Chunk A.2 and Part B rollout).

## Evolution of Project Decisions

- **[Date/Timestamp - e.g., 2025-11-05]**: Initialized Memory Bank.
- **2025-11-05 (Initial Investigation - Signup Error):**
    - Backend `signupHandler` expected wrapped request; OpenAPI & frontend used flat.
    - Initial plan: Align backend to flat request.
- **2025-11-05 (User Directive 1 - Wrapped Requests):**
    - Decision: Enforce `{"data": ...}` wrapper for all API request bodies.
    - Revised plan: Modify OpenAPI & frontend to send wrapped request. Proposed named wrapper schemas.
- **2025-11-05 (User Directive 2 - Inline OpenAPI Wrappers & Testing):**
    - Convention: `{"data": ...}` for requests/success responses; `{"errors": ...}` for errors. Request wrappers in OpenAPI: inline.
    - Plan: Update OpenAPI (inline), then frontend, then backend tests.
- **2025-11-05 (User Directive 3 - Chunked Backend-First Iteration):** Initial chunked plan.
- **2025-11-05 (User Directive 4 - Add Quality Steps):**
    - **Critical Workflow Addition:** `make generate-types` after OpenAPI changes, `make test` after backend changes.
    - **Latest Plan for Signup & Rollout (incorporating quality steps):**
        - **Signup (Part A - Chunked, Backend-First with integrated testing):**
            - A.1: Update OpenAPI (inline request wrapper), run `make generate-types`, verify backend handler, run `make test`. (COMPLETED)
            - A.2: Detailed backend testing (request handling, response shapes), adjust handler response generation if needed, iterating with `make test`. (Partially addressed by validation error fix, ongoing for broader coverage).
            - A.3: Frontend changes for signup payload wrapping. (COMPLETED)
        - **Rollout (Part B - Future):** Apply similar chunked, backend-first, test-driven approach (with `make generate-types` and `make test`) to other endpoints.
- **2025-11-05 (Afternoon - Signup Validation Error Fix):**
    - **Problem:** `TestUserSignup_ValidationErrors` failing due to incorrect API error code (expected `GOSOCIAL-006-VALIDATION_ERROR`, got `GOSOCIAL-001-BAD_REQUEST`) and missing/incorrect field name in error response.
    - **Solution Implemented:**
        1.  `userService.Create` modified to return `validator.ValidationErrors` directly.
        2.  Unit tests for `userService` updated to expect `validator.ValidationErrors`.
        3.  `cmd/api/api.go#writeJSONError` modified to accept a `fieldName` string.
        4.  `cmd/api/api.go#handleErrors` updated:
            - Added case for `validator.ValidationErrors`.
            - Implemented `toSnakeCase` helper to convert struct field names (e.g., "FirstName") to snake_case (e.g., "first_name").
            - Populated `apitypes.ApiError.Field` with the converted field name.
        5.  Calls to `writeJSONError` in `api.go` and `auth_handlers.go` updated.
    - **Outcome:** All unit and functional tests (including `TestUserSignup_ValidationErrors`) passed.
- **2025-11-05 (Afternoon - Frontend Signup Payload Wrapping - Chunk A.3):**
    - **Task:** Update frontend to send wrapped `{"data": ...}` payload for signup.
    - **Changes:**
        1.  Confirmed generated frontend types (`frontend/src/generated/api-types.ts`) correctly define the wrapped request structure for the signup operation.
        2.  Modified `frontend/src/services/authService.ts` in the `signup` method to wrap the `signupData` in a `data` object before making the API call.
    - **Outcome:** User confirmed successful signup via UI and verified correct wrapped request payload. Signup endpoint (Part A of plan) is now fully aligned (OpenAPI, backend, frontend).
- **2025-11-05 (Afternoon - Login Endpoint Request Wrapping - Part B):**
    - **Task:** Apply `{"data": ...}` request payload wrapping convention to `/v1/auth/login`.
    - **Changes:**
        1.  Updated `openapi/v1/paths/auth.yaml` for `/login` requestBody to use inline `data` wrapper.
        2.  Ran `make generate-types`.
        3.  Confirmed backend `loginHandler` in `cmd/api/auth_handlers.go` was already compatible.
        4.  Ran `make test`; all backend tests passed.
        5.  Updated `frontend/src/services/authService.ts#login` method to send the wrapped payload.
    - **Outcome:** User confirmed successful login functionality with the wrapped request.
- **2025-11-05 (Afternoon - Convention Rollout Part B - Backend/Func Tests):**
    - **Task:** Systematically apply `{"data": ...}` request payload wrapping to all relevant V1 POST/PUT endpoints (backend handlers & functional tests).
    - **Endpoints Updated & Verified (Backend & Functional Tests):**
        - `PUT /v1/users`: OpenAPI, handler verified, functional tests updated.
        - `POST /v1/posts` & `PUT /v1/posts/{id}`: OpenAPI, handlers updated, functional tests updated.
        - `POST /v1/posts/{postId}/comments` & `PUT /v1/posts/{postId}/comments/{id}`: OpenAPI, handlers updated, functional tests updated.
    - **Process for each:** Updated OpenAPI, ran `make generate-types`, updated/verified backend handler, updated functional tests, ran `make test` successfully.
    - **Outcome:** All targeted V1 POST/PUT endpoints now adhere to the request wrapping convention on the backend, and their functional tests are aligned.
- **2025-11-05 (Late Afternoon - Frontend Convention Rollout - User Profile):**
    - **Task:** Update frontend for `PUT /v1/users` to send wrapped payload.
    - **Investigation:** Searched `frontend/src` for usage of `/v1/users` (PUT). Found only in generated types.
    - **Outcome:** Determined "update user profile" feature is not yet implemented in frontend. No code changes required for this endpoint.
- **2025-11-05 (Late Afternoon - Frontend Convention Rollout - Posts):**
    - **Task:** Update frontend for `POST /v1/posts` and `PUT /v1/posts/{id}` to send wrapped payloads.
    - **Investigation:** Searched `frontend/src` for usage of `/v1/posts` (POST/PUT). Found only in generated types.
    - **Outcome:** Determined "create/update post" features are not yet implemented in frontend. No code changes required for these endpoints.
- **2025-11-05 (Late Afternoon - Frontend Convention Rollout - Comments):**
    - **Task:** Update frontend for `POST /v1/posts/{postId}/comments` and `PUT /v1/posts/{postId}/comments/{id}` to send wrapped payloads.
    - **Investigation:** Searched `frontend/src` for usage of comment-related paths. Found only in generated types.
    - **Outcome:** Determined "create/update comment" features are not yet implemented in frontend. No code changes required for these endpoints. This concludes the frontend payload convention rollout for existing V1 POST/PUT features.

*(This file will be updated regularly to reflect the project's journey.)*
