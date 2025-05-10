# Active Context

## Current Work Focus

- Updating the Memory Bank with initial project context.
- Preparing to discuss and plan the implementation of a new "likes" endpoint.

## Recent Changes

- Initialized the Memory Bank with placeholder files.
- Updated `projectbrief.md` and `productContext.md` with information from `README.md` and user input.

## Next Steps

- Complete updates for `systemPatterns.md`, `techContext.md`, and `progress.md`.
- Discuss the requirements and design for the new "likes" endpoint.
- Define the OpenAPI specification for the "likes" endpoint.
- Implement the backend logic for "likes".
- Implement the frontend UI and logic for "likes".

## Active Decisions and Considerations

- How to structure the "likes" feature:
    - Can users like posts, comments, or both?
    - What data needs to be stored for a "like" (e.g., user ID, liked item ID, timestamp)?
    - How will "like" counts be efficiently retrieved and displayed?
- Adhering to existing OpenAPI conventions when defining the new endpoint.
- The need for a strategy to manage and version OpenAPI bundled specifications over time (a future task).

## Important Patterns and Preferences

- Use of OpenAPI as the single source of truth for API contracts.
- Code generation for backend (Go) and frontend (TypeScript) types from OpenAPI.
- Facade pattern to decouple application code from generated types.
- Consistent API response structure (`data` for success, `errors` for failure).
- API versioning (currently at v1).
- Use of partial OpenAPI files, bundled using `@redocly/cli`.

## Learnings and Project Insights

- The project has a well-defined structure for API development using OpenAPI.
- Clear separation of concerns between backend and frontend.
- Established processes for dependency management, testing, and running the application (via `Makefile`).
- The `README.md` provides a good overview of the project setup and development workflows.
