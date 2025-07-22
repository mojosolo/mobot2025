// Package validation provides common validation utilities for MoBot 2025
package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	
	"github.com/mojosolo/mobot2025/internal/errors"
)

// File validation

// ValidateFilePath validates a file path
func ValidateFilePath(path string) error {
	if path == "" {
		return errors.New(errors.ErrorTypeValidation, "ValidateFilePath", path, 
			errors.ErrMissingParameter)
	}
	
	// Check for directory traversal
	if strings.Contains(path, "..") {
		return errors.New(errors.ErrorTypeValidation, "ValidateFilePath", path,
			fmt.Errorf("path contains directory traversal"))
	}
	
	return nil
}

// ValidateFileExtension validates file has expected extension
func ValidateFileExtension(path string, validExtensions ...string) error {
	ext := strings.ToLower(filepath.Ext(path))
	
	for _, valid := range validExtensions {
		if ext == strings.ToLower(valid) {
			return nil
		}
	}
	
	return errors.New(errors.ErrorTypeValidation, "ValidateFileExtension", path,
		fmt.Errorf("invalid file extension %s, expected one of %v", ext, validExtensions))
}

// ValidateAEPFile validates an AEP file path
func ValidateAEPFile(path string) error {
	if err := ValidateFilePath(path); err != nil {
		return err
	}
	
	return ValidateFileExtension(path, ".aep", ".aepx")
}

// Data validation

// ValidateStringLength validates string length constraints
func ValidateStringLength(value, name string, min, max int) error {
	length := len(value)
	
	if length < min {
		return errors.New(errors.ErrorTypeValidation, "ValidateStringLength", name,
			fmt.Errorf("length %d is less than minimum %d", length, min))
	}
	
	if max > 0 && length > max {
		return errors.New(errors.ErrorTypeValidation, "ValidateStringLength", name,
			fmt.Errorf("length %d exceeds maximum %d", length, max))
	}
	
	return nil
}

// ValidateRange validates numeric value is within range
func ValidateRange[T int | int32 | int64 | float32 | float64](value T, name string, min, max T) error {
	if value < min {
		return errors.New(errors.ErrorTypeValidation, "ValidateRange", name,
			fmt.Errorf("value %v is less than minimum %v", value, min))
	}
	
	if value > max {
		return errors.New(errors.ErrorTypeValidation, "ValidateRange", name,
			fmt.Errorf("value %v exceeds maximum %v", value, max))
	}
	
	return nil
}

// ValidateRequired validates value is not empty/zero
func ValidateRequired(value interface{}, name string) error {
	// Check various types for empty/zero values
	switch v := value.(type) {
	case string:
		if v == "" {
			return errors.New(errors.ErrorTypeValidation, "ValidateRequired", name,
				errors.ErrMissingParameter)
		}
	case int, int32, int64, float32, float64:
		// Numeric types are always valid (0 might be valid)
		return nil
	case bool:
		// Bool is always valid
		return nil
	case nil:
		return errors.New(errors.ErrorTypeValidation, "ValidateRequired", name,
			errors.ErrMissingParameter)
	default:
		// For other types, assume valid if not nil
		if v == nil {
			return errors.New(errors.ErrorTypeValidation, "ValidateRequired", name,
				errors.ErrMissingParameter)
		}
	}
	
	return nil
}

// Pattern validation

var (
	// Common regex patterns
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	identifierRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	versionRegex = regexp.MustCompile(`^\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$`)
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New(errors.ErrorTypeValidation, "ValidateEmail", email,
			fmt.Errorf("invalid email format"))
	}
	return nil
}

// ValidateIdentifier validates identifier format (starts with letter, alphanumeric + _ -)
func ValidateIdentifier(id string) error {
	if !identifierRegex.MatchString(id) {
		return errors.New(errors.ErrorTypeValidation, "ValidateIdentifier", id,
			fmt.Errorf("invalid identifier format"))
	}
	return nil
}

// ValidateVersion validates semantic version format
func ValidateVersion(version string) error {
	if !versionRegex.MatchString(version) {
		return errors.New(errors.ErrorTypeValidation, "ValidateVersion", version,
			fmt.Errorf("invalid version format, expected X.Y.Z"))
	}
	return nil
}

// Struct validation

// Validator interface for custom validation
type Validator interface {
	Validate() error
}

// ValidateStruct validates a struct that implements Validator
func ValidateStruct(v Validator) error {
	if v == nil {
		return errors.New(errors.ErrorTypeValidation, "ValidateStruct", "",
			errors.ErrMissingParameter)
	}
	
	return v.Validate()
}

// ValidationError represents multiple validation errors
type ValidationError struct {
	Errors []error
}

func (ve *ValidationError) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	
	var msgs []string
	for _, err := range ve.Errors {
		msgs = append(msgs, err.Error())
	}
	
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

// Add adds an error to validation errors
func (ve *ValidationError) Add(err error) {
	if err != nil {
		ve.Errors = append(ve.Errors, err)
	}
}

// HasErrors returns true if there are validation errors
func (ve *ValidationError) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ToError returns error if there are validation errors, nil otherwise
func (ve *ValidationError) ToError() error {
	if ve.HasErrors() {
		return ve
	}
	return nil
}

// NewValidationError creates a new validation error collector
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make([]error, 0),
	}
}

// Batch validation example:
// ve := validation.NewValidationError()
// ve.Add(validation.ValidateRequired(name, "name"))
// ve.Add(validation.ValidateEmail(email))
// ve.Add(validation.ValidateRange(age, "age", 0, 150))
// if err := ve.ToError(); err != nil {
//     return err
// }