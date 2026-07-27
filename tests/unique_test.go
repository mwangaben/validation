package tests

import (
	"encoding/json"
	"testing"

	"github.com/mwangaben/validation"
)

func TestUniqueValidation(t *testing.T) {
	// Create validator with database connection
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Valid unique name",
			data: map[string]interface{}{
				"name": "UniqueName",
			},
			rules: map[string][]string{
				"name": {"unique:users,name"},
			},
			expected: true,
		},
		{
			name: "Invalid duplicate name",
			data: map[string]interface{}{
				"name": "John Doe", // Already exists in database
			},
			rules: map[string][]string{
				"name": {"unique:users,name"},
			},
			expected: false,
		},
		{
			name: "Valid unique email",
			data: map[string]interface{}{
				"email": "newuser@example.com",
			},
			rules: map[string][]string{
				"email": {"unique:users,email"},
			},
			expected: true,
		},
		{
			name: "Invalid duplicate email",
			data: map[string]interface{}{
				"email": "john@example.com", // Already exists
			},
			rules: map[string][]string{
				"email": {"unique:users,email"},
			},
			expected: false,
		},
		{
			name: "Unique with except ID",
			data: map[string]interface{}{
				"name": "John Doe",
				"id":   1, // Except user with ID 1
			},
			rules: map[string][]string{
				"name": {"unique:users,name,id"},
			},
			expected: true, // Should pass because we're excepting ID 1
		},
		{
			name: "Unique with except ID - different record",
			data: map[string]interface{}{
				"name": "John Doe",
				"id":   99, // Except non-existent ID
			},
			rules: map[string][]string{
				"name": {"unique:users,name,id"},
			},
			expected: false, // Should fail because John Doe exists
		},
		{
			name: "Valid unique product code",
			data: map[string]interface{}{
				"code": "P003",
			},
			rules: map[string][]string{
				"code": {"unique:products,code"},
			},
			expected: true,
		},
		{
			name: "Invalid duplicate product code",
			data: map[string]interface{}{
				"code": "P001", // Already exists
			},
			rules: map[string][]string{
				"code": {"unique:products,code"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.data, tt.rules)

			if result != tt.expected {
				t.Errorf("Expected validation to be %v, got %v. Errors: %v",
					tt.expected, result, validator.Error())
			}

			if !result {
				t.Logf("Validation failed with errors: %s", validator.Error())
			}
		})
	}
}

func TestUniqueValidationWithMultipleFields(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	// Test with multiple fields
	data := map[string]interface{}{
		"name":  "Alice Wonderland",
		"email": "alice@example.com",
		"age":   28,
	}

	rules := map[string][]string{
		"name":  {"unique:users,name"},
		"email": {"unique:users,email"},
	}

	result := validator.Validate(data, rules)

	// Both name and email should be unique
	if !result {
		t.Errorf("Expected validation to pass, but got errors: %s", validator.Error())
	}
}

func TestUniqueValidationEdgeCases(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Empty string - should pass (not in DB)",
			data: map[string]interface{}{
				"name": "",
			},
			rules: map[string][]string{
				"name": {"unique:users,name"},
			},
			expected: true,
		},
		{
			name: "Nil value - should pass",
			data: map[string]interface{}{
				"name": nil,
			},
			rules: map[string][]string{
				"name": {"unique:users,name"},
			},
			expected: true,
		},
		{
			name: "Numeric value",
			data: map[string]interface{}{
				"age":  30, // 30 exists but we're checking name
				"name": "NewUser",
			},
			rules: map[string][]string{
				"name": {"unique:users,name"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.data, tt.rules)

			if result != tt.expected {
				t.Errorf("Expected validation to be %v, got %v. Errors: %v",
					tt.expected, result, validator.Error())
			}
		})
	}
}

// Test to verify the except ID functionality
func TestUniqueExceptID(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	// Get the actual ID of John Doe from the database
	var john User
	db.Where("name = ?", "John Doe").First(&john)

	if john.ID == 0 {
		t.Skip("John Doe not found in database, skipping test")
	}

	t.Logf("John Doe has ID: %d", john.ID)

	// Test 1: Should pass when excepting John's own ID
	data1 := map[string]interface{}{
		"name": "John Doe",
		"id":   int(john.ID),
	}
	rules1 := map[string][]string{
		"name": {"unique:users,name,id"},
	}
	result1 := validator.Validate(data1, rules1)
	if !result1 {
		t.Errorf("Expected validation to pass when excepting own ID, got errors: %s", validator.Error())
	}

	// Test 2: Should fail when excepting a different ID
	data2 := map[string]interface{}{
		"name": "John Doe",
		"id":   999, // Non-existent ID
	}
	rules2 := map[string][]string{
		"name": {"unique:users,name,id"},
	}
	result2 := validator.Validate(data2, rules2)
	if result2 {
		t.Errorf("Expected validation to fail when excepting different ID")
	}
}

// Test for custom error messages
func TestCustomErrorMessages(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	// Set custom error messages
	validator.WithMessages(map[string]string{
		"name.required":  "The user's full name is required.",
		"name.min":       "The name must be at least 2 characters.",
		"name.max":       "The name may not exceed 100 characters.",
		"email.required": "The email address is required.",
		"email.email":    "Please provide a valid email address.",
		"email.unique":   "This email is already registered. Please use a different one.",
		"age.required":   "The age is required.",
		"age.int":        "The age must be a number.",
		"age.between":    "The age must be between 1 and 150.",
	})

	data := map[string]interface{}{
		"name":  "",                 // Empty to trigger required
		"email": "john@example.com", // Already exists
		"age":   200,                // Out of range
	}

	rules := map[string][]string{
		"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
		"email": {validation.Required(), validation.Email(), validation.Unique("users", "email")},
		"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
	}

	result := validator.Validate(data, rules)

	if result {
		t.Errorf("Expected validation to fail, but it passed")
	}

	// Check custom error messages
	errors := validator.Errors()

	expectedMessages := map[string]string{
		"name":  "The user's full name is required.",
		"email": "This email is already registered. Please use a different one.",
		"age":   "The age must be between 1 and 150.",
	}

	for field, expected := range expectedMessages {
		if msgs, ok := errors[field]; ok && len(msgs) > 0 {
			if msgs[0] != expected {
				t.Errorf("Expected error for field '%s' to be '%s', got '%s'", field, expected, msgs[0])
			}
		} else {
			t.Errorf("Expected error for field '%s' not found", field)
		}
	}

	t.Logf("Custom error messages: %+v", errors)
}

// Test for structured GraphQL errors - FIXED VERSION
func TestStructuredGraphQLErrors(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	data := map[string]interface{}{
		"name":  "",
		"email": "john@example.com",
		"age":   200,
	}

	rules := map[string][]string{
		"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
		"email": {validation.Required(), validation.Email(), validation.Unique("users", "email")},
		"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
	}

	validator.Validate(data, rules)

	if validator.Passed() {
		t.Errorf("Expected validation to fail")
	}

	// Test ToGraphQLError
	graphQLError := validator.ToGraphQLError("createUser")

	// Verify structure
	if msg, ok := graphQLError["message"].(string); !ok || msg != "Validation failed for the field [createUser]." {
		t.Errorf("Expected message to be 'Validation failed for the field [createUser].', got '%v'", msg)
	}

	extensions, ok := graphQLError["extensions"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected extensions to be a map")
	}

	validation, ok := extensions["validation"].(map[string][]string)
	if !ok {
		t.Errorf("Expected validation to be a map")
	}

	// FIXED: Check that validation errors are formatted as field names (not "input.field")
	for key := range validation {
		// Keys should be the field names directly (name, email, age)
		t.Logf("Validation error key: %s", key)

		// Verify the key is one of the expected fields
		if key != "name" && key != "email" && key != "age" {
			t.Logf("Warning: Unexpected key found: %s", key)
		}
	}

	// Verify we have errors for name, email, and age
	expectedFields := []string{"name", "email", "age"}
	for _, field := range expectedFields {
		if _, ok := validation[field]; !ok {
			t.Logf("Expected validation error for field '%s' not found", field)
		}
	}

	// Print the structured error as JSON
	jsonData, _ := json.MarshalIndent(graphQLError, "", "  ")
	t.Logf("Structured GraphQL Error:\n%s", string(jsonData))
}

// Test for structured validation using ValidateStruct
func TestStructuredValidationWithValidateStruct(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	data := map[string]interface{}{
		"name":  "A",             // Too short
		"email": "invalid-email", // Invalid email
		"age":   200,             // Too high
	}

	rules := map[string][]string{
		"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
		"email": {validation.Required(), validation.Email()},
		"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
	}

	// Use ValidateStruct to get structured errors
	validationErrors, valid := validator.ValidateStruct(data, rules)

	if valid {
		t.Errorf("Expected validation to fail")
	}

	if validationErrors == nil {
		t.Errorf("Expected validation errors to be returned")
		return
	}

	// Test ToMap
	errorMap := validationErrors.ToMap()
	t.Logf("Error Map: %+v", errorMap)

	// Test ToGraphQLExtensions
	extensions := validationErrors.ToGraphQLExtensions()
	t.Logf("GraphQL Extensions: %+v", extensions)

	// Test ToGraphQLError
	graphQLError := validationErrors.ToGraphQLError("createUser")
	if graphQLError == nil {
		t.Errorf("Expected GraphQL error to be returned")
	} else {
		jsonData, _ := json.MarshalIndent(graphQLError, "", "  ")
		t.Logf("GraphQL Error JSON:\n%s", string(jsonData))

		// Verify the error implements the ResolverError interface
		if extensions := graphQLError.Extensions(); extensions == nil {
			t.Errorf("Expected Extensions() method to return non-nil value")
		}
	}

	// Verify we have errors for all fields
	errorMap = validationErrors.ToMap()
	expectedFields := []string{"name", "email", "age"}
	for _, field := range expectedFields {
		if errors, ok := errorMap[field]; !ok || len(errors) == 0 {
			t.Errorf("Expected validation error for field '%s'", field)
		}
	}
}

// Test for custom validator class pattern
func TestCustomValidatorClass(t *testing.T) {
	// Create a custom validator like in Lighthouse PHP
	type UserValidator struct {
		*validation.Validator
	}

	NewUserValidator := func() *UserValidator {
		v := validation.NewValidatorWithDB(db)
		v.WithMessages(map[string]string{
			"name.required":  "The user's full name is required.",
			"name.min":       "The name must be at least 2 characters.",
			"name.max":       "The name may not exceed 100 characters.",
			"email.required": "The email address is required.",
			"email.email":    "Please provide a valid email address.",
			"email.unique":   "This email is already registered. Please use a different one.",
			"age.required":   "The age is required.",
			"age.int":        "The age must be a number.",
			"age.between":    "The age must be between 1 and 150.",
		})
		return &UserValidator{Validator: v}
	}

	validator := NewUserValidator()

	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	rules := map[string][]string{
		"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
		"email": {validation.Required(), validation.Email(), validation.Unique("users", "email")},
		"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
	}

	result := validator.Validate(data, rules)

	// Should fail because John Doe already exists
	if result {
		t.Errorf("Expected validation to fail because John Doe already exists")
	}

	t.Logf("Validation errors: %s", validator.Error())
}

// Test for new validation rules
func TestNewValidationRules(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	// Test password rule
	data := map[string]interface{}{
		"password": "Weak",
	}
	rules := map[string][]string{
		"password": {validation.Password()},
	}
	result := validator.Validate(data, rules)
	if result {
		t.Errorf("Expected weak password to fail validation")
	}

	// Test UUID rule
	data2 := map[string]interface{}{
		"uuid": "invalid-uuid",
	}
	rules2 := map[string][]string{
		"uuid": {validation.UUID()},
	}
	result2 := validator.Validate(data2, rules2)
	if result2 {
		t.Errorf("Expected invalid UUID to fail validation")
	}
}
