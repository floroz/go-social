package errorcodes // Updated package name

// ApiErrorCode defines a type for standard API error codes for type safety.
type ApiErrorCode string

// Standard API Error Codes following GOSOCIAL-<NNN>-<LABEL> format
const (
	CodeBadRequest          ApiErrorCode = "GOSOCIAL-001-BAD_REQUEST"
	CodeUnauthorized        ApiErrorCode = "GOSOCIAL-002-UNAUTHORIZED"
	CodeForbidden           ApiErrorCode = "GOSOCIAL-003-FORBIDDEN"
	CodeNotFound            ApiErrorCode = "GOSOCIAL-004-NOT_FOUND"
	CodeConflict            ApiErrorCode = "GOSOCIAL-005-CONFLICT"
	CodeValidationError     ApiErrorCode = "GOSOCIAL-006-VALIDATION_ERROR"
	CodeInternalServerError ApiErrorCode = "GOSOCIAL-007-INTERNAL_SERVER_ERROR"
)
