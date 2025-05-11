# Active Context

## Current Work Focus

- **Initialization of Memory Bank:** The immediate task is to create the foundational set of Memory Bank documents. This involves `projectbrief.md`, `productContext.md`, `techContext.md`, `systemPatterns.md`, `activeContext.md` (this file), and `progress.md`.

## Recent Changes

- `memory-bank/projectbrief.md` created.
- `memory-bank/productContext.md` created.
- `memory-bank/techContext.md` created.
- `memory-bank/systemPatterns.md` created.

## Next Steps

1. Create `memory-bank/progress.md`.
2. Review all created Memory Bank files for initial completeness and accuracy based on the provided file structure.
3. Await further instructions or tasks from the user.

## Active Decisions and Considerations

- The content of the Memory Bank files is currently based on inferences from the project's file and directory structure. Deeper analysis will require reading specific files (e.g., `go.mod`, `package.json`, `openapi.yaml`, `docker-compose.yaml`, various `.go` and `.ts` files).
- Placeholder names and "To be determined" sections in the Memory Bank files highlight areas needing further investigation or user input.

## Important Patterns and Preferences (from `.clinerules/`)

*   **Generic Functions & `any`**: Use `any` inside generic function bodies if TypeScript cannot match runtime logic to type logic (e.g., conditional return types). Avoid `as <type>` casting generally.
*   **Default Exports**: Avoid default exports unless required by a framework (e.g., Next.js pages). Prefer named exports.
*   **Discriminated Unions**: Proactively use for modeling data with varying shapes (e.g., event types, fetching states) to prevent impossible states. Use `switch` statements for handling.
*   **Enums**: Do not introduce new enums. Retain existing ones. Use `as const` objects for enum-like behavior. Be mindful of numeric enum reverse mapping.
*   **`import type`**: Use `import type` for all type imports, preferably at the top level.
*   **Installing Libraries**: Use package manager commands (`pnpm add`, `yarn add`, `npm install`) to install the latest versions, rather than manually editing `package.json`.
*   **`interface extends`**: Prefer `interface extends` over `&` for modeling inheritance due to performance.
*   **JSDoc Comments**: Use for functions and types if behavior isn't self-evident. Use `@link` for internal references.
*   **Naming Conventions**:
    *   kebab-case for files (`my-component.ts`)
    *   camelCase for variables/functions (`myVariable`)
    *   UpperCamelCase (PascalCase) for classes/types/interfaces (`MyClass`)
    *   ALL_CAPS for constants/enum values (`MAX_COUNT`)
    *   `T` prefix for generic type parameters (`TKey`)
*   **`noUncheckedIndexedAccess`**: Be aware that if enabled, object/array indexing returns `T | undefined`.
*   **Optional Properties**: Use sparingly. Prefer `property: T | undefined` over `property?: T` if the property's presence is critical but its value can be absent.
*   **`readonly` Properties**: Use by default for object types to prevent accidental mutation. Omit only if genuinely mutable.
*   **Return Types**: Declare return types for top-level module functions (except JSX components).
*   **Throwing Errors**: Consider result types (`Result<T, E>`) for operations that might fail (e.g., parsing JSON) instead of throwing, unless throwing is idiomatic for the framework (e.g., backend request handlers).

## Learnings and Project Insights

- The project "Go Social" is a full-stack application with a Go backend and a React/TypeScript frontend.
- It utilizes OpenAPI for API design and code generation.
- Docker is used for containerization.
- A comprehensive set of `.clinerules` dictates coding standards and best practices for TypeScript development.

*(This file will be updated frequently as work progresses.)*
