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

### Proposed OpenAPI Specification for "Likes" (for discussion)

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
          type: string
          description: Unique identifier for the like.
          readOnly: true
          example: "lk_clxkz2kv0000008jy1g7g3h7n"
        user_id:
          type: string # Corresponds to User ID. If User.id is int64, this should be int64 too.
          description: ID of the user who liked the content.
          readOnly: true # Set by the backend based on authenticated user.
          example: "usr_clxkyv02o000008jye1g2f3h5" # Example if User ID is string
          # if User ID is int64:
          # type: integer
          # format: int64
          # example: 101
        post_id:
          type: string # Corresponds to Post ID. If Post.id is int64, this should be int64.
          description: ID of the post being liked (mutually exclusive with comment_id).
          nullable: true
          example: "post_clxkz0q9k000008jya1b7g2e3" # Example if Post ID is string
          # if Post ID is int64:
          # type: integer
          # format: int64
          # example: 201
        comment_id:
          type: string # Corresponds to Comment ID. If Comment.id is int64, this should be int64.
          description: ID of the comment being liked (mutually exclusive with post_id).
          nullable: true
          example: "cmt_clxkz1b2a000008jyf4g6h7i8" # Example if Comment ID is string
          # if Comment ID is int64:
          # type: integer
          # format: int64
          # example: 301
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
          type: string # Should match Post.id type (string or integer)
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
                $ref: '#/components/schemas/ListLikesSuccessResponse'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '500':
          $ref: '#/components/responses/InternalServerError'
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
                $ref: '#/components/schemas/CreateLikeSuccessResponse'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '409':
          $ref: '#/components/responses/ConflictError'
        '500':
          $ref: '#/components/responses/InternalServerError'
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
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /v1/comments/{commentId}/likes:
    parameters:
      - name: commentId
        in: path
        required: true
        description: The ID of the comment.
        schema:
          type: string # Should match Comment.id type (string or integer)
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
                $ref: '#/components/schemas/ListLikesSuccessResponse'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '500':
          $ref: '#/components/responses/InternalServerError'
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
                $ref: '#/components/schemas/CreateLikeSuccessResponse'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '409':
          $ref: '#/components/responses/ConflictError'
        '500':
          $ref: '#/components/responses/InternalServerError'
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
          $ref: '#/components/responses/UnauthorizedError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '500':
          $ref: '#/components/responses/InternalServerError'
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
  #   UnauthorizedError: ...
  #   NotFoundError: ...
  #   ConflictError: ...
  #   InternalServerError: ...
# ...
```

**Discussion Points for this Proposal:**
1.  **ID Types:** Confirm if `Like` related IDs (`id`, `user_id`, `post_id`, `comment_id`) should be `string` (as per Go struct) or `integer` (like other entities).
2.  **Listing Likes Content:** Is an array of full `Like` objects suitable for `GET` requests, or would a list of user IDs/count be preferred in some cases?
3.  **Common Error Responses:** Ensure `#/components/responses/*` refs point to correctly defined shared responses.

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
