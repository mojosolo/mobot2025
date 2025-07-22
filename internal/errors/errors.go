// Package errors provides common error types and utilities for MoBot 2025
package errors

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	// File errors
	ErrFileNotFound      = errors.New("file not found")
	ErrInvalidFormat     = errors.New("invalid file format")
	ErrUnsupportedFormat = errors.New("unsupported file format")
	ErrCorruptedFile     = errors.New("file appears to be corrupted")
	
	// Parsing errors
	ErrInvalidHeader     = errors.New("invalid file header")
	ErrUnexpectedEOF     = errors.New("unexpected end of file")
	ErrInvalidChunk      = errors.New("invalid chunk format")
	ErrMissingRequiredData = errors.New("missing required data")
	
	// Database errors
	ErrDatabaseLocked    = errors.New("database is locked")
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateRecord   = errors.New("duplicate record")
	ErrTransactionFailed = errors.New("transaction failed")
	
	// Validation errors
	ErrValidationFailed  = errors.New("validation failed")
	ErrInvalidInput      = errors.New("invalid input")
	ErrMissingParameter  = errors.New("missing required parameter")
	ErrOutOfRange        = errors.New("value out of range")
	
	// Workflow errors
	ErrWorkflowFailed    = errors.New("workflow execution failed")
	ErrTaskFailed        = errors.New("task execution failed")
	ErrDependencyFailed  = errors.New("dependency not met")
	ErrTimeout           = errors.New("operation timed out")
)

// ErrorType represents categories of errors
type ErrorType string

const (
	ErrorTypeFile       ErrorType = "file"
	ErrorTypeParse      ErrorType = "parse"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeWorkflow   ErrorType = "workflow"
	ErrorTypeInternal   ErrorType = "internal"
)

// Error represents a detailed error with context
type Error struct {
	Type    ErrorType
	Op      string      // Operation that failed
	Path    string      // File path or resource identifier
	Err     error       // Underlying error
	Context interface{} // Additional context
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s %s: %s: %v", e.Type, e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("%s %s: %v", e.Type, e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Is implements errors.Is
func (e *Error) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// New creates a new detailed error
func New(errType ErrorType, op, path string, err error) *Error {
	return &Error{
		Type: errType,
		Op:   op,
		Path: path,
		Err:  err,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, op string) error {
	if err == nil {
		return nil
	}
	
	// If it's already our error type, add operation context
	if e, ok := err.(*Error); ok {
		e.Op = fmt.Sprintf("%s->%s", op, e.Op)
		return e
	}
	
	// Otherwise create new error
	return &Error{
		Type: ErrorTypeInternal,
		Op:   op,
		Err:  err,
	}
}

// IsFileError checks if error is file-related
func IsFileError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeFile
	}
	return errors.Is(err, ErrFileNotFound) || 
	       errors.Is(err, ErrInvalidFormat) ||
	       errors.Is(err, ErrUnsupportedFormat) ||
	       errors.Is(err, ErrCorruptedFile)
}

// IsParseError checks if error is parsing-related
func IsParseError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeParse
	}
	return errors.Is(err, ErrInvalidHeader) ||
	       errors.Is(err, ErrUnexpectedEOF) ||
	       errors.Is(err, ErrInvalidChunk) ||
	       errors.Is(err, ErrMissingRequiredData)
}

// IsDatabaseError checks if error is database-related
func IsDatabaseError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeDatabase
	}
	return errors.Is(err, ErrDatabaseLocked) ||
	       errors.Is(err, ErrRecordNotFound) ||
	       errors.Is(err, ErrDuplicateRecord) ||
	       errors.Is(err, ErrTransactionFailed)
}

// IsValidationError checks if error is validation-related
func IsValidationError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeValidation
	}
	return errors.Is(err, ErrValidationFailed) ||
	       errors.Is(err, ErrInvalidInput) ||
	       errors.Is(err, ErrMissingParameter) ||
	       errors.Is(err, ErrOutOfRange)
}

// IsRetryable checks if error can be retried
func IsRetryable(err error) bool {
	return errors.Is(err, ErrDatabaseLocked) ||
	       errors.Is(err, ErrTimeout) ||
	       errors.Is(err, ErrTransactionFailed)
}