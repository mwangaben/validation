package main

import (
	"fmt"
	"github.com/mwangaben/validation"
)

func main() {
	validator := validation.NewValidator()

	// Test data with various scenarios
	data := map[string]interface{}{
		"username":   "johndoe123",
		"email":      "john@example.com",
		"password":   "SecurePass123",
		"age":        25,
		"website":    "https://example.com",
		"role":       "admin",
		"score":      95.5,
		"is_active":  true,
		"phone":      "+1234567890",
		"uuid":       "123e4567-e89b-12d3-a456-426614174000",
		"birth_date": "1990-01-01",
		"tags":       []string{"go", "validation"},
	}

	// Rules as map[string][]string
	rules := map[string][]string{
		"username":   {"required", "string", "min:3", "max:50", "alpha_num"},
		"email":      {"required", "email"},
		"password":   {"required", "string", "min:8", "password"},
		"age":        {"required", "int", "between:18,100"},
		"website":    {"url"},
		"role":       {"required", "in:admin,user,guest"},
		"score":      {"numeric", "between:0,100"},
		"is_active":  {"boolean"},
		"phone":      {"phone"},
		"uuid":       {"uuid"},
		"birth_date": {"date"},
		"tags":       {"array"},
	}

	fmt.Println("🔍 Validating data...")
	if validator.Validate(data, rules) {
		fmt.Println("✅ All validations passed!")
	} else {
		fmt.Println("❌ Validation failed:")
		fmt.Println(validator.Error())

		// You can also get detailed errors
		fmt.Println("\n📋 Detailed errors:")
		for field, errors := range validator.Errors() {
			fmt.Printf("  %s:\n", field)
			for _, err := range errors {
				fmt.Printf("    - %s\n", err)
			}
		}
	}
}
