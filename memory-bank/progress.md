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
- **Current Focus:** Revising the plan to fix the signup 400 error by adopting a new project convention: all API request bodies (like success responses) should be wrapped with a `{"data": ...}` structure.
- **Blockers:** None currently for planning. Awaiting user approval for the revised plan.

## Known Issues

- **Signup Endpoint 400 Error (json: unknown field "first_name"):**
    - **Status:** Investigated. Root cause understood. Plan revised based on new project convention.
    - **Original Description:** Backend `signupHandler` expected a `{"data": ...}` wrapped request, but OpenAPI and frontend used a flat request.
    - **Revised Understanding & Path Forward:** The project will now enforce `{"data": ...}` wrappers for *all* request bodies. The fix involves:
        1. Modifying the OpenAPI spec for `/v1/auth/signup` (and potentially others) to define wrapped request bodies.
        2. Regenerating frontend types.
        3. Updating the frontend to send the wrapped payload.
        4. The backend `signupHandler`'s current expectation of a wrapped request will then be correct.
    - **Impact:** Users cannot sign up via the frontend until this is resolved.
- None identified yet.

## Evolution of Project Decisions

- **[Date/Timestamp - e.g., 2025-11-05]**: Initialized Memory Bank.
- **2025-11-05 (Initial Investigation - Signup Error):**
    - **Finding:** Backend `signupHandler` expected a `{"data": ...}` wrapped request, while OpenAPI and frontend used a flat request. Field name `first_name` (snake_case) was consistent.
    - **Initial Plan:** Modify backend handler to accept flat request, aligning with then-current OpenAPI.
- **2025-11-05 (Revised Strategy - User Directive):**
    - **New Convention:** Decided to enforce `{"data": ...}` wrappers for all API *request bodies* for consistency with success response wrappers.
    - **Revised Plan for Signup Error:**
        - Modify OpenAPI specification for `/v1/auth/signup` request body to be wrapped (e.g., `{"data": {"first_name": ...}}`).
        - Regenerate frontend types.
        - Update frontend to send the wrapped payload.
        - The backend `signupHandler`'s existing expectation for a wrapped request becomes the correct state.
    - This shifts the "source of truth" for request structure to the new convention, requiring OpenAPI and frontend updates.

*(This file will be updated regularly to reflect the project's journey.)*
