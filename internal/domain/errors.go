package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type ErrorDetail struct {
	// TODO: in the future an error catalog can be used to map errors to codes
	// Code    int    `json:"code"`
	Message string `json:"message"`
}

type BadRequestError struct {
	ErrorDetail
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return &BadRequestError{
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}

type ValidationError struct {
	ErrorDetail
	Field string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field: field,
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}

type InternalServerError struct {
	ErrorDetail
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) error {
	return &InternalServerError{
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}
