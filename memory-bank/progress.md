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

- **Phase:** Debugging / Maintenance.
- **Current Focus:** Investigating and planning a fix for a 400 error on the signup endpoint (`/v1/auth/signup`). Investigation complete, technical plan pending presentation.
- **Blockers:** None currently for planning. Awaiting user approval to implement the fix.

## Known Issues

- **Signup Endpoint 400 Error (json: unknown field "first_name"):**
    - **Status:** Investigated. Root cause identified.
    - **Description:** The backend's `signupHandler` in `cmd/api/auth_handlers.go` incorrectly expects the signup request payload to be nested under a `data` key (e.g., `{"data": {"first_name": "..."}}`).
    - **Conflict:** This contradicts the OpenAPI specification and the frontend implementation, both of which correctly use a flat payload structure (e.g., `{"first_name": "..."}`).
    - **Impact:** Users cannot sign up via the frontend.
- None identified yet.

## Evolution of Project Decisions

- **[Date/Timestamp - e.g., 2025-11-05]**: Initialized Memory Bank. Content is based on inferences from file structure and provided `.clinerules`. Further code/configuration file analysis is needed to confirm and elaborate on these initial assessments.
- **2025-11-05**: Investigated a 400 error (`json: unknown field "first_name"`) on the `/v1/auth/signup` endpoint.
    - **Finding:** The error is due to the backend handler (`cmd/api/auth_handlers.go`) expecting a JSON payload nested under a `data` key, while the OpenAPI specification (and thus the frontend and its generated types) defines a flat payload.
    - **Field Name Consistency:** The actual field name `first_name` (snake_case) is consistent between the OpenAPI spec, frontend-generated types, and the backend's `domain.CreateUserDTO` struct tags. The issue is purely the payload structure (nesting).

*(This file will be updated regularly to reflect the project's journey.)*
