# Active Context

## Current Work Focus

- Refining project scope in Memory Bank and planning the OpenAPI specification for a new 'likes' endpoint.

## Recent Changes

- Initialized the Memory Bank with placeholder files.
- Updated `projectbrief.md`, `productContext.md`, `systemPatterns.md`, and `techContext.md` with information from `README.md` and user input.

## Next Steps

1.  Finalize scope adjustments in Memory Bank files (`activeContext.md`, `projectbrief.md`, `progress.md`).
2.  Discuss and define the requirements and design for the new "likes" endpoint (e.g., what can be liked, data structure for a like).
3.  Create/update the necessary OpenAPI partial files (e.g., new schemas in `openapi/shared/schemas/` or `openapi/v1/schemas/`, new paths in `openapi/v1/paths/`).
4.  Update `openapi/openapi.yaml` to reference the new "likes" paths and schemas.
5.  Run `make generate-types` (or at least the `redocly bundle` part) to create an updated `openapi/openapi-bundled.yaml`.
6.  Confirm the new "likes" endpoint is correctly represented in the Swagger UI.

## Active Decisions and Considerations

- How to structure the "likes" feature within the OpenAPI specification:
    - Can users like posts, comments, or both? This will determine the path structure and request/response schemas.
    - What data needs to be represented in the "like" schema (e.g., user ID, liked item ID, timestamp)?
- Adhering to existing OpenAPI conventions (e.g., response wrappers, error handling, naming) when defining the new endpoint.
- The implementation of backend/frontend logic for "likes" is out of scope for the current task.
- The need for a strategy to manage and version OpenAPI bundled specifications over time remains a future task.

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
