package apitypes

import (
	"github.com/floroz/go-social/internal/generated"
	"github.com/oapi-codegen/runtime/types"
)

// Re-export generated types with potentially cleaner names or just as aliases.
// This acts as a facade over the raw generated types.

// Schemas
type LoginRequest = generated.LoginRequest
type LoginResponse = generated.LoginResponse
type SignupRequest = generated.SignupRequest
type User = generated.User

type ApiError = generated.ApiError
type ApiErrorResponse = generated.ApiErrorResponse
type SignupSuccessResponse = generated.SignupSuccessResponse
type LoginSuccessResponse = generated.LoginSuccessResponse

// User endpoint types
type UpdateUserProfileRequest = generated.UpdateUserProfileRequest
type GetUserProfileSuccessResponse = generated.GetUserProfileSuccessResponse
type UpdateUserProfileSuccessResponse = generated.UpdateUserProfileSuccessResponse

// Post endpoint types
type Post = generated.Post // Shared Post schema
type CreatePostRequest = generated.CreatePostRequest
type UpdatePostRequest = generated.UpdatePostRequest
type CreatePostSuccessResponse = generated.CreatePostSuccessResponse
type GetPostSuccessResponse = generated.GetPostSuccessResponse
type UpdatePostSuccessResponse = generated.UpdatePostSuccessResponse
type ListPostsSuccessResponse = generated.ListPostsSuccessResponse

// Comment endpoint types
type Comment = generated.Comment // Shared Comment schema
type CreateCommentRequest = generated.CreateCommentRequest
type UpdateCommentRequest = generated.UpdateCommentRequest
type CreateCommentSuccessResponse = generated.CreateCommentSuccessResponse
type GetCommentSuccessResponse = generated.GetCommentSuccessResponse
type UpdateCommentSuccessResponse = generated.UpdateCommentSuccessResponse
type ListCommentsSuccessResponse = generated.ListCommentsSuccessResponse

// Runtime Types (if needed directly, like Email)
type Email = types.Email

// Add other schemas as they are defined and generated...
// Example:
// type Post = generated.Post
// type Comment = generated.Comment
// type CreatePostRequest = generated.CreatePostRequest

// Note: We are not re-exporting request body types like
// LoginUserV1JSONRequestBody as those are specific to the generator's client/server code
// which we are not using directly in this way. We use the schema types (LoginRequest etc.)
// in our handlers.
