# Active Context

## Current Work Focus

- Finalizing Memory Bank updates after successfully defining the OpenAPI specification for the "likes" endpoint.

## Recent Changes

- Initialized the Memory Bank with placeholder files.
- Updated `projectbrief.md`, `productContext.md`, `systemPatterns.md`, and `techContext.md` with information from `README.md` and user input.
- Iteratively defined and corrected the OpenAPI specification for the "likes" feature:
    - Created `openapi/v1/schemas/like.yaml` (with integer IDs).
    - Created `openapi/v1/paths/like.yaml` (with relative path refs for schemas).
    - Updated `openapi/openapi.yaml` with the "Likes V1" tag and correct path/schema references.
- Successfully ran `make generate-types`, bundling the spec and generating types.
- Confirmed (by user) that the new "likes" endpoints are correctly represented in Swagger UI.
- Updated Memory Bank (`activeContext.md`, `progress.md`) to reflect troubleshooting steps and then successful completion.

## Next Steps

1.  **Finalize Memory Bank updates to reflect task completion (this step).**
2.  (No further technical steps for this task)

## Active Decisions and Considerations

- How to structure the "likes" feature within the OpenAPI specification:
    - Can users like posts, comments, or both? This will determine the path structure and request/response schemas.
    - What data needs to be represented in the "like" schema (e.g., user ID, liked item ID, timestamp)?
- Adhering to existing OpenAPI conventions (e.g., response wrappers, error handling, naming) when defining the new endpoint.
- The implementation of backend/frontend logic for "likes" is out of scope for the current task.
- The need for a strategy to manage and version OpenAPI bundled specifications over time remains a future task.

### Finalized OpenAPI Specification for "Likes" (as implemented)

**1. New File: `openapi/v1/schemas/like.yaml`**
```yaml
# openapi/v1/schemas/like.yaml
components:
  schemas:
    Like:
      type: object
      description: Represents a user's like on a post or comment. One of post_id or comment_id should be populated.
      properties:
        id:
          type: integer
          format: int64
          description: Unique identifier for the like.
          readOnly: true
          example: 1001
        user_id:
          type: integer
          format: int64
          description: ID of the user who liked the content.
          readOnly: true # Set by the backend based on authenticated user.
          example: 101
        post_id:
          type: integer
          format: int64
          description: ID of the post being liked (mutually exclusive with comment_id).
          nullable: true
          example: 201
        comment_id:
          type: integer
          format: int64
          description: ID of the comment being liked (mutually exclusive with post_id).
          nullable: true
          example: 301
        created_at:
          type: string
          format: date-time
          description: Timestamp when the like was created.
          readOnly: true
          example: '2024-03-10T10:30:00Z'
      required:
        - id
        - user_id
        - created_at

    CreateLikeSuccessResponse:
      type: object
      description: Standard wrapper for the successful like creation response.
      properties:
        data:
          $ref: '#/components/schemas/Like'
      required:
        - data

    ListLikesSuccessResponse:
      type: object
      description: Standard wrapper for listing likes.
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/Like'
      required:
        - data
```

**2. New File: `openapi/v1/paths/like.yaml`**
```yaml
# openapi/v1/paths/like.yaml
paths:
  /v1/posts/{postId}/likes:
    parameters:
      - name: postId
        in: path
        required: true
        description: The ID of the post.
        schema:
          type: integer
          format: int64
    get:
      tags:
        - Likes V1
      summary: List likes for a post
      description: Retrieves a list of likes for a specific post.
      operationId: listPostLikesV1
      security:
        - bearerAuth: []
      responses:
        '200':
          description: A list of likes retrieved successfully.
          content:
            application/json:
              schema:
                $ref: '../schemas/like.yaml#/components/schemas/ListLikesSuccessResponse'
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Post not found.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error retrieving likes.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
    post:
      tags:
        - Likes V1
      summary: Like a post
      description: Creates a new like on a specific post for the authenticated user.
      operationId: createPostLikeV1
      security:
        - bearerAuth: []
      responses:
        '201':
          description: Post liked successfully. Returns the created like object.
          content:
            application/json:
              schema:
                $ref: '../schemas/like.yaml#/components/schemas/CreateLikeSuccessResponse'
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Post not found.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '409':
          description: Post already liked by the user or other conflict.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error creating like.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
    delete:
      tags:
        - Likes V1
      summary: Unlike a post
      description: Removes the authenticated user's like from a specific post.
      operationId: deletePostLikeV1
      security:
        - bearerAuth: []
      responses:
        '204':
          description: Post unliked successfully. No content returned.
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Post not found or like not found for the user.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error deleting like.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'

  /v1/comments/{commentId}/likes:
    parameters:
      - name: commentId
        in: path
        required: true
        description: The ID of the comment.
        schema:
          type: integer
          format: int64
    get:
      tags:
        - Likes V1
      summary: List likes for a comment
      description: Retrieves a list of likes for a specific comment.
      operationId: listCommentLikesV1
      security:
        - bearerAuth: []
      responses:
        '200':
          description: A list of likes retrieved successfully.
          content:
            application/json:
              schema:
                $ref: '../schemas/like.yaml#/components/schemas/ListLikesSuccessResponse'
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Comment not found.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error retrieving likes.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
    post:
      tags:
        - Likes V1
      summary: Like a comment
      description: Creates a new like on a specific comment for the authenticated user.
      operationId: createCommentLikeV1
      security:
        - bearerAuth: []
      responses:
        '201':
          description: Comment liked successfully. Returns the created like object.
          content:
            application/json:
              schema:
                $ref: '../schemas/like.yaml#/components/schemas/CreateLikeSuccessResponse'
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Comment not found.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '409':
          description: Comment already liked by the user or other conflict.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error creating like.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
    delete:
      tags:
        - Likes V1
      summary: Unlike a comment
      description: Removes the authenticated user's like from a specific comment.
      operationId: deleteCommentLikeV1
      security:
        - bearerAuth: []
      responses:
        '204':
          description: Comment unliked successfully. No content returned.
        '401':
          description: Authentication required or invalid token.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '404':
          description: Comment not found or like not found for the user.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
        '500':
          description: Server error deleting like.
          content:
            application/json:
              schema:
                $ref: '../../shared/schemas/common.yaml#/components/schemas/ApiErrorResponse'
```

**3. Updates to `openapi/openapi.yaml` (relevant parts)**
```yaml
# openapi/openapi.yaml
# ...
tags:
  # ... existing tags ...
  - name: Likes V1 # New Tag
    description: Operations related to likes on posts and comments (Version 1)

paths:
  # ... existing path references ...
  # New Path References for Likes
  /v1/posts/{postId}/likes:
    $ref: './v1/paths/like.yaml#/paths/~1v1~1posts~1{postId}~1likes'
  /v1/comments/{commentId}/likes:
    $ref: './v1/paths/like.yaml#/paths/~1v1~1comments~1{commentId}~1likes'

components:
  schemas:
    # ... existing schema references ...
    # New Schema References for Likes
    Like:
      $ref: './v1/schemas/like.yaml#/components/schemas/Like'
    CreateLikeSuccessResponse:
      $ref: './v1/schemas/like.yaml#/components/schemas/CreateLikeSuccessResponse'
    ListLikesSuccessResponse:
      $ref: './v1/schemas/like.yaml#/components/schemas/ListLikesSuccessResponse'
  
  # responses: # Common responses should be defined centrally
  # responses:
  #   Common response components like UnauthorizedError, NotFoundError etc. are not explicitly defined here.
  #   Instead, individual operations directly reference #/components/schemas/ApiErrorResponse for their error states,
  #   with the expectation that the backend will populate the 'code' and 'message' fields appropriately for each error type.
  #   If more specific, reusable error response components are desired later, they can be added here or in shared/schemas/common.yaml.
# ...
```

**Summary of Changes Based on Feedback:**
1.  **ID Types:** All `Like` related IDs (`id`, `user_id`, `post_id`, `comment_id`) and path parameters (`postId`, `commentId`) are now `type: integer`, `format: int64`. Examples updated.
2.  **Listing Likes Content:** Remains as an array of full `Like` objects.
3.  **Common Error Responses:** Error responses now directly reference `#/components/schemas/ApiErrorResponse`. Descriptions have been added to each error status code for clarity.

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
