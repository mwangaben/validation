package validation

import (
	"encoding/json"
	"fmt"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string   `json:"field"`
	Message string   `json:"message"`
	Rules   []string `json:"rules,omitempty"`
}

// ValidationErrors collection of validation errors
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors(errors map[string][]string) *ValidationErrors {
	return &ValidationErrors{
		Errors: errors,
	}
}

// Error implements error interface
func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Errors)
}

// ToJSON returns JSON representation
func (e *ValidationErrors) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// HasErrors checks if there are any errors
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// Get returns errors for a specific field
func (e *ValidationErrors) Get(field string) []string {
	if errors, ok := e.Errors[field]; ok {
		return errors
	}
	return []string{}
}

// First returns the first error message for a field
func (e *ValidationErrors) First(field string) string {
	if errors, ok := e.Errors[field]; ok && len(errors) > 0 {
		return errors[0]
	}
	return ""
}

// All returns all error messages
func (e *ValidationErrors) All() map[string][]string {
	return e.Errors
}

// ToMap converts errors to a map
func (e *ValidationErrors) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	for field, messages := range e.Errors {
		if len(messages) == 1 {
			result[field] = messages[0]
		} else {
			result[field] = messages
		}
	}
	return result
}
