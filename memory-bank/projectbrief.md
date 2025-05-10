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
- Implement new features such as "likes" for posts/comments.
- Develop a strategy for producing and persisting different OpenAPI bundle snapshots.

## Scope

- Backend development in Go (Chi Router, PostgreSQL).
- Frontend development in React (TypeScript, Vite, Tailwind CSS, Shadcn/ui, Zustand, React Query).
- Database migrations using golang-migrate.
- API definition and code generation based on OpenAPI specs.
- Initial focus on implementing a "likes" feature.
- Future consideration for OpenAPI bundle snapshotting.
