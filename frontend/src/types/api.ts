// This file acts as a facade for the raw generated API types.
// Import necessary types from the generated file and re-export them,
// potentially with cleaner names or grouped logically.

import { components } from "../generated/api-types";

// --- Schemas ---
export type User = components["schemas"]["User"];
export type LoginRequest = components["schemas"]["LoginRequest"];
export type SignupRequest = components["schemas"]["SignupRequest"];

export type ApiError = components["schemas"]["ApiError"];
export type ApiErrorResponse = components["schemas"]["ApiErrorResponse"];
export type SignupSuccessResponse =
  components["schemas"]["SignupSuccessResponse"];
export type LoginSuccessResponse =
  components["schemas"]["LoginSuccessResponse"];
export type LoginResponse = components["schemas"]["LoginResponse"];

// Add other schema exports as needed
// export type Post = components["schemas"]["Post"];
// export type Comment = components["schemas"]["Comment"];

// --- Operations ---
// We might not need to re-export operations directly if services/hooks
// handle the request/response types based on the schemas above.
// If needed, they could be exported like:
// export type SignupUserOperation = operations["signupUserV1"];
// export type LoginUserOperation = operations["loginUserV1"];

// --- Utility Types (Optional) ---
// Example: Extracting the request body type for a specific operation
// export type SignupRequestBody =
//   operations["signupUserV1"]["requestBody"]["content"]["application/json"];
