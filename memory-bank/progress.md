# Project Progress

## What Works

- **Initial Project Structure:** The basic directory structure for a Go backend and a React/TypeScript frontend is in place.
- **Memory Bank Initialization:** The core Memory Bank files have been created.
- **API Error Handling (Signup - Validation & Conflicts):**
    - User signup validation errors (e.g., password too short, invalid email format) now correctly return HTTP 400 with API error code `GOSOCIAL-006-VALIDATION_ERROR` and include the specific problematic field name (e.g., "password", "email", "first_name") in the response.
    - User signup attempts with duplicate email/username now correctly return HTTP 409 with API error code `GOSOCIAL-005-CONFLICT`.
    - Functional test `TestUserSignup_ValidationErrors` and related tests for duplicate signups are passing.

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

- **Phase:** API Implementation & Testing.
- **Current Focus:** Finalizing Memory Bank update after successfully fixing `TestUserSignup_ValidationErrors` functional test.
- **Blockers:** None for the current memory bank update task.

## Known Issues

- **Signup Endpoint Payload Wrapping:**
    - **Status:** Backend expects wrapped (`{"data": ...}`) requests for signup. OpenAPI was updated (Chunk A.1) to reflect this with an inline wrapper. Frontend changes to send wrapped requests are pending (was part of a previous Chunk A.3).
    - **Impact:** Frontend signup might still fail if it sends a flat payload, until frontend is updated.
    - **Note:** This is distinct from the recently fixed validation error code/field name issue.
- **Comprehensive Backend Testing for Response Shapes:** While the signup validation error response shape is now correct, a broader review and addition of tests for other endpoints' request/response shapes (ensuring `data` and `errors` wrappers are consistently tested) is an ongoing effort.

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
            - A.3 (Later): Frontend changes for payload wrapping.
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

*(This file will be updated regularly to reflect the project's journey.)*
