# Progress

## What Works

- Core application structure (Go backend, React frontend).
- User authentication (signup, login, logout, token refresh).
- Post creation, retrieval, updating, and deletion.
- Comment creation, retrieval, updating, and deletion on posts.
- Database setup with PostgreSQL and migrations using golang-migrate.
- API definition using OpenAPI, with code generation for Go and TypeScript types.
- Development environment setup with Docker, Makefiles for common tasks.
- Frontend setup with Vite, TypeScript, Tailwind, Zustand, React Query.
- Basic testing setup for backend and frontend.
- Interactive API documentation via Swagger UI.

## What's Left to Build

- **New Features:**
    - "Likes" functionality for posts and/or comments. This is the immediate next feature.
- **Technical Improvements/Future Considerations:**
    - Strategy for producing and persisting different OpenAPI bundle snapshots over time.
    - Potentially more comprehensive testing.
    - Further UI/UX enhancements on the frontend.
    - Scalability and performance optimizations as the user base grows.

## Current Status

- The Memory Bank has been initialized and populated with initial project context based on `README.md` and user input.
- The next immediate task is to discuss, plan, and implement the "likes" feature.

## Known Issues

- No specific issues identified from the provided context, apart from the "missing feature" of OpenAPI bundle snapshotting.

## Evolution of Project Decisions

- **Initial Decision:** Adopt OpenAPI for API design to ensure consistency and enable code generation.
- **Initial Decision:** Structure the project with a Go backend and React frontend.
- **Initial Decision:** Use PostgreSQL as the database.
- **Initial Decision:** Implement API versioning from the start (e.g., `/v1/`).
- **Current State:** The project has followed these initial decisions and established a solid foundation. The focus is now shifting towards adding new features like "likes" and considering longer-term improvements like OpenAPI snapshot management.
