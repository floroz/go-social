package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return &BadRequestError{
		Message: message,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewConflictError(message string) error {
	return &ConflictError{
		Message: message,
	}
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) error {
	return &InternalServerError{
		Message: message,
	}
}
