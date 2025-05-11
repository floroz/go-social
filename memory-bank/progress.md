# Project Progress

## What Works

- **Initial Project Structure:** The basic directory structure for a Go backend and a React/TypeScript frontend is in place.
- **Memory Bank Initialization:** The core Memory Bank files (`projectbrief.md`, `productContext.md`, `techContext.md`, `systemPatterns.md`, `activeContext.md`, `progress.md`) have been created with initial content based on file structure analysis and user-provided rules.

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

- **Phase:** API Design Refinement & Debugging.
- **Current Focus:** Finalizing a plan to fix the signup 400 error. This involves adopting a project-wide convention for `{"data": ...}` wrapped API request bodies (defined inline in OpenAPI), ensuring backend tests cover response shapes, and outlining a rollout strategy for other endpoints.
- **Blockers:** None currently for planning. Awaiting user approval for the latest revised plan.

## Known Issues

- **Signup Endpoint 400 Error (json: unknown field "first_name"):**
    - **Status:** Investigated. Root cause understood. Detailed plan formulated based on new project conventions.
    - **Path Forward:**
        1. Modify OpenAPI spec for `/v1/auth/signup` to define an *inline* `{"data": ...}` wrapper for the request body.
        2. Regenerate frontend types.
        3. Update frontend to send the wrapped payload.
        4. Verify backend `signupHandler` (which expects wrapped request) functions correctly.
        5. Add/update backend tests for request unmarshaling and response structure (ensuring `{"data":...}` for success, `{"errors":...}` for errors).
    - **Impact:** Users cannot sign up via the frontend until this is resolved.
- **Potential Lack of Backend Tests for Response Shapes:** As part of the signup fix, it's necessary to verify/add tests that explicitly check the JSON structure of API responses (e.g., presence of `data` or `errors` wrappers). This might be a gap across other endpoints too.

## Evolution of Project Decisions

- **[Date/Timestamp - e.g., 2025-11-05]**: Initialized Memory Bank.
- **2025-11-05 (Initial Investigation - Signup Error):**
    - Backend `signupHandler` expected wrapped request; OpenAPI & frontend used flat.
    - Initial plan: Align backend to flat request (matching then-current OpenAPI).
- **2025-11-05 (User Directive 1 - Wrapped Requests):**
    - Decision: Enforce `{"data": ...}` wrapper for all API request bodies.
    - Revised plan: Modify OpenAPI & frontend to send wrapped request; backend handler's expectation becomes correct. Proposed using named wrapper schemas in OpenAPI.
- **2025-11-05 (User Directive 2 - Inline OpenAPI Wrappers & Testing):**
    - **New Convention Details:**
        - `{"data": ...}` wrappers for requests and success responses; `{"errors": ...}` for error responses.
        - Request body wrappers in OpenAPI to be defined *inline* in path definitions, not as separate named schemas.
    - **Latest Plan for Signup & Rollout:**
        - **Signup:** Implement inline `data` wrapper in OpenAPI for signup request, update frontend, verify backend, add comprehensive backend tests for request/response shapes.
        - **Rollout:** Iteratively apply this pattern (inline OpenAPI request wrappers, frontend changes, backend handler updates if needed, testing) to other POST/PUT/PATCH endpoints.

*(This file will be updated regularly to reflect the project's journey.)*
