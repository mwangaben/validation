package tests

import (
	"testing"

	"github.com/mwangaben/validation"
)

func TestExistsValidation(t *testing.T) {
	// Create validator with database connection
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Valid existing name",
			data: map[string]interface{}{
				"name": "John Doe", // Exists in database
			},
			rules: map[string][]string{
				"name": {"exists:users,name"},
			},
			expected: true,
		},
		{
			name: "Invalid non-existing name",
			data: map[string]interface{}{
				"name": "NonExistentUser",
			},
			rules: map[string][]string{
				"name": {"exists:users,name"},
			},
			expected: false,
		},
		{
			name: "Valid existing email",
			data: map[string]interface{}{
				"email": "john@example.com", // Exists in database
			},
			rules: map[string][]string{
				"email": {"exists:users,email"},
			},
			expected: true,
		},
		{
			name: "Invalid non-existing email",
			data: map[string]interface{}{
				"email": "nonexistent@example.com",
			},
			rules: map[string][]string{
				"email": {"exists:users,email"},
			},
			expected: false,
		},
		{
			name: "Valid existing product code",
			data: map[string]interface{}{
				"code": "P001", // Exists in database
			},
			rules: map[string][]string{
				"code": {"exists:products,code"},
			},
			expected: true,
		},
		{
			name: "Invalid non-existing product code",
			data: map[string]interface{}{
				"code": "P999",
			},
			rules: map[string][]string{
				"code": {"exists:products,code"},
			},
			expected: false,
		},
		{
			name: "Valid existing ID",
			data: map[string]interface{}{
				"id": 1, // Exists in database
			},
			rules: map[string][]string{
				"id": {"exists:users,id"},
			},
			expected: true,
		},
		{
			name: "Invalid non-existing ID",
			data: map[string]interface{}{
				"id": 999,
			},
			rules: map[string][]string{
				"id": {"exists:users,id"},
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

func TestExistsWithMultipleConditions(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	// Test with multiple fields
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	rules := map[string][]string{
		"name":  {"exists:users,name"},
		"email": {"exists:users,email"},
	}

	result := validator.Validate(data, rules)

	// Both name and email should exist
	if !result {
		t.Errorf("Expected validation to pass, but got errors: %s", validator.Error())
	}
}

func TestExistsEdgeCases(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Empty string - should fail (empty doesn't exist)",
			data: map[string]interface{}{
				"name": "",
			},
			rules: map[string][]string{
				"name": {"exists:users,name"},
			},
			expected: false,
		},
		{
			name: "Nil value - should fail",
			data: map[string]interface{}{
				"name": nil,
			},
			rules: map[string][]string{
				"name": {"exists:users,name"},
			},
			expected: false,
		},
		{
			name: "Numeric value as string",
			data: map[string]interface{}{
				"age": "30",
			},
			rules: map[string][]string{
				"age": {"exists:users,age"},
			},
			expected: true, // 30 exists in the database
		},
		{
			name: "Numeric value not existing",
			data: map[string]interface{}{
				"age": 99,
			},
			rules: map[string][]string{
				"age": {"exists:users,age"},
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
		})
	}
}

// TestExistsWithDifferentTables tests the exists rule across different tables
func TestExistsWithDifferentTables(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Check in users table",
			data: map[string]interface{}{
				"name": "Jane Smith",
			},
			rules: map[string][]string{
				"name": {"exists:users,name"},
			},
			expected: true,
		},
		{
			name: "Check in products table",
			data: map[string]interface{}{
				"code": "P002",
			},
			rules: map[string][]string{
				"code": {"exists:products,code"},
			},
			expected: true,
		},
		{
			name: "Cross-table check - name in products table",
			data: map[string]interface{}{
				"name": "John Doe",
			},
			rules: map[string][]string{
				"name": {"exists:products,name"},
			},
			expected: false, // John Doe is in users, not products
		},
		{
			name: "Cross-table check - code in users table",
			data: map[string]interface{}{
				"code": "P001",
			},
			rules: map[string][]string{
				"code": {"exists:users,code"},
			},
			expected: false, // P001 is in products, not users
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

// TestExistsCombinedWithOtherRules tests exists with other validation rules
func TestExistsCombinedWithOtherRules(t *testing.T) {
	validator := validation.NewValidatorWithDB(db)

	tests := []struct {
		name     string
		data     map[string]interface{}
		rules    map[string][]string
		expected bool
	}{
		{
			name: "Exists and email validation",
			data: map[string]interface{}{
				"email": "john@example.com",
			},
			rules: map[string][]string{
				"email": {"required", "email", "exists:users,email"},
			},
			expected: true,
		},
		{
			name: "Exists and string validation",
			data: map[string]interface{}{
				"name": "John Doe",
			},
			rules: map[string][]string{
				"name": {"required", "string", "min:2", "exists:users,name"},
			},
			expected: true,
		},
		{
			name: "Exists fails with invalid email",
			data: map[string]interface{}{
				"email": "notanemail",
			},
			rules: map[string][]string{
				"email": {"required", "email", "exists:users,email"},
			},
			expected: false, // Fails email validation first
		},
		{
			name: "Exists with non-existing value",
			data: map[string]interface{}{
				"name": "NonExistentUser",
			},
			rules: map[string][]string{
				"name": {"required", "string", "min:2", "exists:users,name"},
			},
			expected: false, // Fails exists validation
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

// Benchmark tests for exists rule
func BenchmarkExistsValidation(b *testing.B) {
	validator := validation.NewValidatorWithDB(db)

	data := map[string]interface{}{
		"name": "John Doe",
	}

	rules := map[string][]string{
		"name": {"exists:users,name"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(data, rules)
	}
}

// TestExistsWithCaseInsensitiveCheck (optional - depends on database collation)
func TestExistsWithCaseInsensitiveCheck(t *testing.T) {
	// This test may vary depending on your database collation
	validator := validation.NewValidatorWithDB(db)

	// In MySQL with utf8mb4_unicode_ci, this should be case-insensitive
	data := map[string]interface{}{
		"name": "john doe", // lowercase version of "John Doe"
	}

	rules := map[string][]string{
		"name": {"exists:users,name"},
	}

	result := validator.Validate(data, rules)

	// This might pass or fail depending on your MySQL collation
	// utf8mb4_unicode_ci is case-insensitive, so it should find "John Doe"
	if !result {
		t.Logf("Case-insensitive check failed. Check your database collation.")
		t.Logf("Error: %s", validator.Error())
	} else {
		t.Log("✅ Case-insensitive check passed")
	}
}
