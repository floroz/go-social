# Project Brief

## Core Requirements

- Social media application with Go backend and React frontend.
- Features for users to create posts, comment on posts, and interact.
- Follows OpenAPI specification for API design.
- API versioning to avoid breaking changes.

## Goals

- Provide a robust and scalable social media platform.
- Maintain a clear and consistent API contract using OpenAPI.
- Enable future expansion for client applications beyond the current frontend.
- Define OpenAPI specifications for new features, starting with "likes".
- Implement new features (backend and frontend) based on defined OpenAPI specs (future work).
- Develop a strategy for producing and persisting different OpenAPI bundle snapshots (future work).

## Scope

- **Current Task:** Define the OpenAPI specification for a "likes" feature. This includes:
    - Creating or updating relevant YAML partials for paths and schemas.
    - Updating the main `openapi.yaml` file.
    - Generating the `openapi-bundled.yaml` and verifying its correctness in Swagger UI.
- **Out of Scope for Current Task:**
    - Implementation of backend logic for the "likes" feature.
    - Implementation of frontend UI/logic for the "likes" feature.
    - Developing the strategy for OpenAPI bundle snapshotting.
- **Overall Project Scope (includes future work):**
    - Backend development in Go (Chi Router, PostgreSQL).
    - Frontend development in React (TypeScript, Vite, Tailwind CSS, Shadcn/ui, Zustand, React Query).
    - Database migrations using golang-migrate.
    - Full API lifecycle management based on OpenAPI specs (definition, code generation, implementation).
