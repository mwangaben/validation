package validation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationError represents a single validation error with field information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation passed"
	}

	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Field+": "+err.Message)
	}
	return strings.Join(messages, "; ")
}

// ToMap converts validation errors to a map format (field => []messages)
func (ve *ValidationErrors) ToMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range ve.Errors {
		field := err.Field
		if idx := strings.Index(field, "."); idx != -1 {
			field = field[:idx]
		}
		message := strings.TrimPrefix(err.Message, "The ")
		message = strings.TrimSuffix(message, ".")
		message = strings.TrimSpace(message)
		result[field] = append(result[field], message)
	}
	return result
}

// ToGraphQLExtensions converts to GraphQL extension format
func (ve *ValidationErrors) ToGraphQLExtensions() map[string]interface{} {
	if len(ve.Errors) == 0 {
		return nil
	}
	return map[string]interface{}{
		"validation": ve.ToMap(),
	}
}

// ToGraphQLError creates a GraphQL error with extensions
// This is the main method to use in resolvers
func (ve *ValidationErrors) ToGraphQLError(operationName string) *ValidationGraphQLError {
	if len(ve.Errors) == 0 {
		return nil
	}
	return &ValidationGraphQLError{
		Message:        fmt.Sprintf("Validation failed for the field [%s].", operationName),
		Path:           []interface{}{operationName},
		ExtensionsData: ve.ToGraphQLExtensions(),
	}
}

// ToJSON returns the validation errors as JSON
func (ve *ValidationErrors) ToJSON() string {
	bytes, _ := json.Marshal(ve.ToMap())
	return string(bytes)
}

// HasErrors checks if there are any validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// GetFieldErrors returns errors for a specific field
func (ve *ValidationErrors) GetFieldErrors(field string) []string {
	var result []string
	for _, err := range ve.Errors {
		if err.Field == field {
			result = append(result, err.Message)
		}
	}
	return result
}

// ValidationGraphQLError is a custom error type for GraphQL validation errors
// This implements the ResolverError interface for graph-gophers/graphql-go
type ValidationGraphQLError struct {
	Message        string                 `json:"message"`
	Path           []interface{}          `json:"path,omitempty"`
	ExtensionsData map[string]interface{} `json:"extensions"`
}

// Error implements the error interface
func (e *ValidationGraphQLError) Error() string {
	return e.Message
}

// Extensions implements the ResolverError interface for graph-gophers/graphql-go
func (e *ValidationGraphQLError) Extensions() map[string]interface{} {
	return e.ExtensionsData
}
