package validation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Validator handles validation
type Validator struct {
	errors   map[string][]string
	data     map[string]interface{}
	rules    map[string][]string
	db       interface{}
	messages map[string]string // Custom error messages
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors:   make(map[string][]string),
		messages: make(map[string]string),
	}
}

// NewValidatorWithDB creates a new validator instance with database connection
func NewValidatorWithDB(db interface{}) *Validator {
	return &Validator{
		errors:   make(map[string][]string),
		db:       db,
		messages: make(map[string]string),
	}
}

// WithMessages sets custom error messages
func (v *Validator) WithMessages(messages map[string]string) *Validator {
	v.messages = messages
	return v
}

// SetMessage sets a custom error message for a specific rule
func (v *Validator) SetMessage(rule, message string) *Validator {
	v.messages[rule] = message
	return v
}

// Validate validates data against rules
func (v *Validator) Validate(data map[string]interface{}, rules map[string][]string) bool {
	v.data = data
	v.rules = rules
	v.errors = make(map[string][]string)

	for field, fieldRules := range rules {
		value := data[field]
		for _, rule := range fieldRules {
			if !v.validateRule(field, value, rule, data) {
				v.addError(field, v.getErrorMessage(field, rule))
			}
		}
	}

	return len(v.errors) == 0
}

// ValidateStruct performs validation and returns structured errors
func (v *Validator) ValidateStruct(data map[string]interface{}, rules map[string][]string) (*ValidationErrors, bool) {
	valid := v.Validate(data, rules)
	if valid {
		return nil, true
	}
	return v.ParseValidationErrors(), false
}

// ParseValidationErrors converts the raw error messages to structured format
func (v *Validator) ParseValidationErrors() *ValidationErrors {
	errors := &ValidationErrors{}

	if len(v.errors) == 0 {
		return errors
	}

	for field, messages := range v.errors {
		for _, message := range messages {
			cleanMessage := strings.TrimPrefix(message, "The ")
			cleanMessage = strings.TrimSuffix(cleanMessage, ".")
			cleanMessage = strings.TrimSpace(cleanMessage)

			if strings.HasPrefix(cleanMessage, field+": ") {
				cleanMessage = strings.TrimPrefix(cleanMessage, field+": ")
			}

			fieldName := field
			if idx := strings.Index(field, "."); idx != -1 {
				fieldName = field[:idx]
			}

			errors.Errors = append(errors.Errors, ValidationError{
				Field:   fieldName,
				Message: cleanMessage,
			})
		}
	}

	return errors
}

// validateRule validates a single rule
func (v *Validator) validateRule(field string, value interface{}, rule string, data map[string]interface{}) bool {
	ruleParts := strings.Split(rule, ":")
	ruleName := ruleParts[0]
	ruleParams := []string{}
	if len(ruleParts) > 1 {
		ruleParams = strings.Split(ruleParts[1], ",")
	}

	switch ruleName {
	case "required":
		return v.validateRequired(value)
	case "email":
		return v.validateEmail(value)
	case "min":
		return v.validateMin(value, ruleParams[0])
	case "max":
		return v.validateMax(value, ruleParams[0])
	case "string":
		return v.validateString(value)
	case "int":
		return v.validateInt(value)
	case "numeric":
		return v.validateNumeric(value)
	case "unique":
		return v.validateUnique(value, data, ruleParams[0], ruleParams[1:]...)
	case "exists":
		return v.validateExists(value, ruleParams[0], ruleParams[1:]...)
	case "in":
		return v.validateIn(value, ruleParams)
	case "not_in":
		return v.validateNotIn(value, ruleParams)
	case "confirmed":
		return v.validateConfirmed(field, value)
	case "date":
		return v.validateDate(value)
	case "url":
		return v.validateURL(value)
	case "alpha":
		return v.validateAlpha(value)
	case "alpha_num":
		return v.validateAlphaNum(value)
	case "boolean":
		return v.validateBoolean(value)
	case "array":
		return v.validateArray(value)
	case "between":
		if len(ruleParams) == 2 {
			return v.validateBetween(value, ruleParams[0], ruleParams[1])
		}
		return true
	case "phone":
		return v.validatePhone(value)
	case "password":
		return v.validatePassword(value)
	case "uuid":
		return v.validateUUID(value)
	default:
		return true
	}
}

// addError adds an error message
func (v *Validator) addError(field, message string) {
	if _, exists := v.errors[field]; !exists {
		v.errors[field] = []string{}
	}
	v.errors[field] = append(v.errors[field], message)
}

// getErrorMessage returns the error message for a field
func (v *Validator) getErrorMessage(field, rule string) string {
	ruleParts := strings.Split(rule, ":")
	ruleName := ruleParts[0]
	param := ""
	if len(ruleParts) > 1 {
		param = ruleParts[1]
	}

	// Check for custom message for this specific field and rule
	customKey := fmt.Sprintf("%s.%s", field, ruleName)
	if msg, ok := v.messages[customKey]; ok {
		return msg
	}

	// Check for custom message for this rule
	if msg, ok := v.messages[ruleName]; ok {
		return msg
	}

	// Default messages
	messages := map[string]string{
		"required":  fmt.Sprintf("The %s field is required.", field),
		"email":     fmt.Sprintf("The %s must be a valid email address.", field),
		"string":    fmt.Sprintf("The %s must be a string.", field),
		"int":       fmt.Sprintf("The %s must be an integer.", field),
		"numeric":   fmt.Sprintf("The %s must be numeric.", field),
		"min":       fmt.Sprintf("The %s must be at least %s.", field, param),
		"max":       fmt.Sprintf("The %s may not be greater than %s.", field, param),
		"in":        fmt.Sprintf("The selected %s is invalid.", field),
		"not_in":    fmt.Sprintf("The selected %s is invalid.", field),
		"confirmed": fmt.Sprintf("The %s confirmation does not match.", field),
		"unique":    fmt.Sprintf("The %s has already been taken.", field),
		"exists":    fmt.Sprintf("The selected %s is invalid.", field),
		"date":      fmt.Sprintf("The %s is not a valid date.", field),
		"url":       fmt.Sprintf("The %s format is invalid.", field),
		"alpha":     fmt.Sprintf("The %s may only contain letters.", field),
		"alpha_num": fmt.Sprintf("The %s may only contain letters and numbers.", field),
		"boolean":   fmt.Sprintf("The %s field must be true or false.", field),
		"array":     fmt.Sprintf("The %s must be an array.", field),
		"between":   fmt.Sprintf("The %s must be between %s.", field, param),
		"phone":     fmt.Sprintf("The %s must be a valid phone number.", field),
		"password":  fmt.Sprintf("The %s must be at least 8 characters with at least one uppercase, one lowercase, and one number.", field),
		"uuid":      fmt.Sprintf("The %s must be a valid UUID.", field),
	}

	if msg, ok := messages[ruleName]; ok {
		return msg
	}
	return fmt.Sprintf("The %s field is invalid.", field)
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string][]string {
	return v.errors
}

// Error returns a formatted error string
func (v *Validator) Error() string {
	var messages []string
	for field, fieldErrors := range v.errors {
		for _, err := range fieldErrors {
			messages = append(messages, fmt.Sprintf("%s: %s", field, err))
		}
	}
	return strings.Join(messages, "; ")
}

// ToJSON returns JSON representation of errors
func (v *Validator) ToJSON(operationName string) string {
	errData := v.ToGraphQLError(operationName)
	data, _ := json.MarshalIndent(errData, "", "  ")
	return string(data)
}

// Passed checks if validation passed
func (v *Validator) Passed() bool {
	return len(v.errors) == 0
}

// Failed checks if validation failed
func (v *Validator) Failed() bool {
	return len(v.errors) > 0
}

// GetValidationErrors returns validation errors as a map
func (v *Validator) GetValidationErrors() map[string][]string {
	return v.errors
}

// HasError checks if a specific field has errors
func (v *Validator) HasError(field string) bool {
	_, exists := v.errors[field]
	return exists
}

// GetFieldErrors returns errors for a specific field
func (v *Validator) GetFieldErrors(field string) []string {
	if errors, ok := v.errors[field]; ok {
		return errors
	}
	return []string{}
}

// ToGraphQLError returns a validation error that can be returned as an error
func (v *Validator) ToGraphQLError(operationName string) map[string]interface{} {
	validationErrors := make(map[string][]string)
	for field, messages := range v.errors {
		fieldName := field
		if idx := strings.Index(field, "."); idx != -1 {
			fieldName = field[:idx]
		}
		validationErrors[fieldName] = append(validationErrors[fieldName], messages...)
	}

	return map[string]interface{}{
		"message": fmt.Sprintf("Validation failed for the field [%s].", operationName),
		"path":    []interface{}{operationName},
		"extensions": map[string]interface{}{
			"validation": validationErrors,
		},
	}
}

// ToGraphQLErrorString returns a JSON string of the GraphQL error
func (v *Validator) ToGraphQLErrorString(operationName string) string {
	errData := v.ToGraphQLError(operationName)
	data, _ := json.Marshal(errData)
	return string(data)
}

// GraphQLError returns an error that implements the error interface
func (v *Validator) GraphQLError(operationName string) error {
	return &ValidationGraphQLError{
		Message: fmt.Sprintf("Validation failed for the field [%s].", operationName),
		Path:    []interface{}{operationName},
		ExtensionsData: map[string]interface{}{
			"validation": v.ToGraphQLError(operationName)["extensions"].(map[string]interface{})["validation"],
		},
	}
}

// GraphQLErrorWithExtensions returns an error with GraphQL extensions
func (v *Validator) GraphQLErrorWithExtensions(operationName string) *ValidationGraphQLError {
	validationErrors := make(map[string][]string)
	for field, messages := range v.errors {
		fieldName := field
		if idx := strings.Index(field, "."); idx != -1 {
			fieldName = field[:idx]
		}
		validationErrors[fieldName] = append(validationErrors[fieldName], messages...)
	}

	return &ValidationGraphQLError{
		Message: fmt.Sprintf("Validation failed for the field [%s].", operationName),
		Path:    []interface{}{operationName},
		ExtensionsData: map[string]interface{}{
			"validation": validationErrors,
		},
	}
}

// GetStructuredErrors returns structured validation errors
func (v *Validator) GetStructuredErrors() *ValidationErrors {
	return v.ParseValidationErrors()
}
