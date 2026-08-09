package tests

import (
	"fmt"
	"github.com/mwangaben/validation/tests/helpers"
	"testing"

	"github.com/mwangaben/validation"
)

func BenchmarkExistsValidation(b *testing.B) {

	var err error
	testDB, err = helpers.InitDB()
	if err != nil {
		_ = fmt.Errorf("failed to connect: %v", err)
	}
	validator := validation.NewValidatorWithDB(testDB)

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
	validator := validation.NewValidatorWithDB(testDB)

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
