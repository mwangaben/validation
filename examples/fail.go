package main

import (
	"fmt"
	"github.com/mwangaben/validation"
)

func main() {
	validator := validation.NewValidator()

	// Invalid data
	data := map[string]interface{}{
		"name":  "A",          // Too short (min:2)
		"email": "notanemail", // Invalid email
		"age":   200,          // Too high (between:1,150)
	}

	rules := map[string][]string{
		"name":  {"required", "string", "min:2", "max:100"},
		"email": {"required", "email"},
		"age":   {"required", "int", "between:1,150"},
	}

	fmt.Println("🔍 Validating invalid data...")
	if validator.Validate(data, rules) {
		fmt.Println("✅ Validation passed")
	} else {
		fmt.Println("❌ Validation failed with errors:")
		fmt.Println(validator.Error())

		// Check individual fields
		if validator.HasError("name") {
			fmt.Printf("Name errors: %v\n", validator.GetFieldErrors("name"))
		}
		if validator.HasError("email") {
			fmt.Printf("Email errors: %v\n", validator.GetFieldErrors("email"))
		}
		if validator.HasError("age") {
			fmt.Printf("Age errors: %v\n", validator.GetFieldErrors("age"))
		}
	}
}
