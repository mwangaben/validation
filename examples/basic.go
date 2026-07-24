package main

import (
	"fmt"
	"github.com/mwangaben/validation"
)

func main() {
	validator := validation.NewValidator()

	// Data to validate
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	// Define rules as map[string][]string
	rules := map[string][]string{
		"name":  {"required", "string", "min:2", "max:100"},
		"email": {"required", "email"},
		"age":   {"required", "int", "between:1,150"},
	}

	// Validate the data
	if validator.Validate(data, rules) {
		fmt.Println("✅ Validation passed")
	} else {
		fmt.Println("❌ Validation failed")
		fmt.Println("Errors:", validator.Error())
	}
}
