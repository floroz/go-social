package apitypes

import (
	"github.com/floroz/go-social/internal/generated"
	"github.com/oapi-codegen/runtime/types"
)

// Re-export generated types with potentially cleaner names or just as aliases.
// This acts as a facade over the raw generated types.

// Schemas
type ErrorResponse = generated.ErrorResponse
type LoginRequest = generated.LoginRequest
type LoginResponse = generated.LoginResponse
type SignupRequest = generated.SignupRequest
type User = generated.User

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
