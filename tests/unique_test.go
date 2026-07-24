package tests

import (
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
