package utils

import (
	"errors"
	"net/http"
	"strings"
)

// CustomError represents a custom error with status code
type CustomError struct {
	StatusCode int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new custom error
func NewCustomError(statusCode int, message string) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// NewBadRequestError creates a 400 error
func NewBadRequestError(message string) *CustomError {
	return NewCustomError(http.StatusBadRequest, message)
}

// NewUnauthorizedError creates a 401 error
func NewUnauthorizedError(message string) *CustomError {
	return NewCustomError(http.StatusUnauthorized, message)
}

// NewForbiddenError creates a 403 error
func NewForbiddenError(message string) *CustomError {
	return NewCustomError(http.StatusForbidden, message)
}

// NewNotFoundError creates a 404 error
func NewNotFoundError(message string) *CustomError {
	return NewCustomError(http.StatusNotFound, message)
}

// NewConflictError creates a 409 error
func NewConflictError(message string) *CustomError {
	return NewCustomError(http.StatusConflict, message)
}

// NewInternalServerError creates a 500 error
func NewInternalServerError(message string) *CustomError {
	return NewCustomError(http.StatusInternalServerError, message)
}

// IsDuplicateError checks if the error is a duplicate entry error
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := err.Error()
	return strings.Contains(errorStr, "Duplicate entry") ||
		strings.Contains(errorStr, "23505") || // PostgreSQL unique violation
		strings.Contains(errorStr, "1062")     // MySQL duplicate entry
}

// IsForeignKeyError checks if the error is a foreign key constraint error
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := err.Error()
	return strings.Contains(errorStr, "foreign key constraint") ||
		strings.Contains(errorStr, "1452") || // MySQL foreign key constraint
		strings.Contains(errorStr, "23503")   // PostgreSQL foreign key violation
}

// IsRecordNotFoundError checks if the error is a record not found error
func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, errors.New("record not found")) ||
		strings.Contains(err.Error(), "record not found")
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]string, 0),
	}
}

// AddError adds an error to the validation result
func (vr *ValidationResult) AddError(error string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, error)
}

// HasErrors checks if there are any validation errors
func (vr *ValidationResult) HasErrors() bool {
	return !vr.IsValid
}