package tests

import (
	"testing"

	"github.com/mwangaben/validation"
)

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

func BenchmarkUniqueValidation(b *testing.B) {
	validator := validation.NewValidatorWithDB(db)

	data := map[string]interface{}{
		"name": "UniqueName",
	}

	rules := map[string][]string{
		"name": {"unique:users,name"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(data, rules)
	}
}
