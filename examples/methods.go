package main

import (
	"fmt"
	"github.com/mwangaben/validation"
)

func main() {
	validator := validation.NewValidator()

	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	rules := map[string][]string{
		"name": {"required", "string", "min:2"},
		"age":  {"required", "int"},
	}

	validator.Validate(data, rules)

	// Check validation status
	fmt.Printf("Passed: %v\n", validator.Passed())
	fmt.Printf("Failed: %v\n", validator.Failed())

	// Get all errors
	fmt.Printf("Errors: %v\n", validator.Errors())

	// Check specific field
	if validator.HasError("name") {
		fmt.Printf("Name has errors: %v\n", validator.GetFieldErrors("name"))
	}
}
