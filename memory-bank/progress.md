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
- **Current Focus:** Preparing to implement the first chunk (Chunk A.1) of the plan to fix the signup 400 error. This involves updating the OpenAPI spec for the signup request with an inline `data` wrapper and then verifying the backend handler structure.
- **Blockers:** Awaiting user approval for the chunked plan.

## Known Issues

- **Signup Endpoint 400 Error (json: unknown field "first_name"):**
    - **Status:** Investigated. Root cause understood. Detailed, chunked, backend-first iterative plan formulated.
    - **Path Forward (Chunked):**
        - **Chunk A.1 (OpenAPI & Backend Handler Verification):** Update OpenAPI for signup request (inline `data` wrapper). Verify backend handler's existing request unmarshaling structure. Confirm OpenAPI for success response wrapper.
        - **Chunk A.2 (Backend Testing & Adjustments):** Run backend tests. Update/add tests for wrapped request handling and, crucially, for `data`/`errors` wrapped response structures. Adjust handler's response generation if needed.
        - **Chunk A.3 (Frontend - Later):** Regenerate types, update frontend to send wrapped request, test.
    - **Impact:** Users cannot sign up via the frontend until this is resolved.
- **Potential Lack of Backend Tests for Response Shapes:** A key part of Chunk A.2 is to address this for the signup endpoint. This needs to be a focus during the rollout to other endpoints as well.

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
- **2025-11-05 (User Directive 3 - Chunked Backend-First Iteration):**
    - **Latest Plan for Signup & Rollout:**
        - **Signup (Part A - Chunked, Backend-First):**
            - A.1: Update OpenAPI (inline request wrapper), verify backend handler structure.
            - A.2: Backend testing (request handling, response shapes), adjust handler response generation if needed.
            - A.3 (Later): Frontend changes.
        - **Rollout (Part B - Future):** Apply similar chunked, backend-first, test-driven approach to other endpoints.

*(This file will be updated regularly to reflect the project's journey.)*
