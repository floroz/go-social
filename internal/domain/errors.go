package domain

import (
	"errors"
	"fmt"
	"net/http"
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
	StatusCode int
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return &BadRequestError{
		StatusCode: http.StatusBadRequest,
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
	StatusCode int
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) error {
	return &InternalServerError{
		StatusCode: http.StatusInternalServerError,
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}

type NotFoundError struct {
	ErrorDetail
	StatusCode int
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {
	return &NotFoundError{
		StatusCode: http.StatusNotFound,
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}

type UnauthorizedError struct {
	ErrorDetail
	StatusCode int
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) error {
	return &UnauthorizedError{
		StatusCode: http.StatusUnauthorized,
		ErrorDetail: ErrorDetail{
			Message: message,
		},
	}
}
