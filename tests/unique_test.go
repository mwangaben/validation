package tests

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mwangaben/validation"
)

var _ = Describe("Unique Validation", func() {
	var validator *validation.Validator

	BeforeEach(func() {
		validator = validation.NewValidatorWithDB(db)
	})

	Context("Single field unique validation", func() {
		It("should validate unique name", func() {
			data := map[string]interface{}{
				"name": "UniqueName",
			}
			rules := map[string][]string{
				"name": {"unique:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject duplicate name", func() {
			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"unique:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("already been taken"))
		})

		It("should validate unique email", func() {
			data := map[string]interface{}{
				"email": "newuser@example.com",
			}
			rules := map[string][]string{
				"email": {"unique:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject duplicate email", func() {
			data := map[string]interface{}{
				"email": "john@example.com",
			}
			rules := map[string][]string{
				"email": {"unique:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should validate unique product code", func() {
			data := map[string]interface{}{
				"code": "P003",
			}
			rules := map[string][]string{
				"code": {"unique:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject duplicate product code", func() {
			data := map[string]interface{}{
				"code": "P001",
			}
			rules := map[string][]string{
				"code": {"unique:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Unique with except ID", func() {
		It("should pass when excepting own ID", func() {
			data := map[string]interface{}{
				"name": "John Doe",
				"id":   1,
			}
			rules := map[string][]string{
				"name": {"unique:users,name,id"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should fail when excepting different ID", func() {
			data := map[string]interface{}{
				"name": "John Doe",
				"id":   99,
			}
			rules := map[string][]string{
				"name": {"unique:users,name,id"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Multiple fields unique validation", func() {
		It("should validate multiple unique fields", func() {
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
			Expect(result).To(BeTrue())
		})
	})

	Context("Unique validation edge cases", func() {
		It("should pass with empty string", func() {
			data := map[string]interface{}{
				"name": "",
			}
			rules := map[string][]string{
				"name": {"unique:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should pass with nil value", func() {
			data := map[string]interface{}{
				"name": nil,
			}
			rules := map[string][]string{
				"name": {"unique:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})
	})

	Context("Custom error messages", func() {
		It("should use custom error messages", func() {
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
				"name":  "",
				"email": "john@example.com",
				"age":   200,
			}
			rules := map[string][]string{
				"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
				"email": {validation.Required(), validation.Email(), validation.Unique("users", "email")},
				"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
			}

			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())

			errors := validator.Errors()
			Expect(errors["name"][0]).To(Equal("The user's full name is required."))
			Expect(errors["email"][0]).To(Equal("This email is already registered. Please use a different one."))
			Expect(errors["age"][0]).To(Equal("The age must be between 1 and 150."))
		})
	})

	Context("Structured GraphQL errors", func() {
		It("should return structured GraphQL errors", func() {
			validator = validation.NewValidatorWithDB(db)

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
			Expect(validator.Passed()).To(BeFalse())

			graphQLError := validator.ToGraphQLError("createUser")
			Expect(graphQLError["message"]).To(Equal("Validation failed for the field [createUser]."))

			extensions, ok := graphQLError["extensions"].(map[string]interface{})
			Expect(ok).To(BeTrue())

			validation, ok := extensions["validation"].(map[string][]string)
			Expect(ok).To(BeTrue())

			expectedFields := []string{"name", "email", "age"}
			for _, field := range expectedFields {
				Expect(validation).To(HaveKey(field))
			}

			jsonData, _ := json.MarshalIndent(graphQLError, "", "  ")
			GinkgoWriter.Printf("Structured GraphQL Error:\n%s", string(jsonData))
		})
	})

	Context("ValidateStruct method", func() {
		It("should return structured validation errors", func() {
			data := map[string]interface{}{
				"name":  "A",
				"email": "invalid-email",
				"age":   200,
			}
			rules := map[string][]string{
				"name":  {validation.Required(), validation.String(), validation.Min(2), validation.Max(100)},
				"email": {validation.Required(), validation.Email()},
				"age":   {validation.Required(), validation.Int(), validation.Between(1, 150)},
			}

			validationErrors, valid := validator.ValidateStruct(data, rules)
			Expect(valid).To(BeFalse())
			Expect(validationErrors).ToNot(BeNil())

			errorMap := validationErrors.ToMap()
			GinkgoWriter.Printf("Error Map: %+v", errorMap)

			graphQLError := validationErrors.ToGraphQLError("createUser")
			Expect(graphQLError).ToNot(BeNil())

			jsonData, _ := json.MarshalIndent(graphQLError, "", "  ")
			GinkgoWriter.Printf("GraphQL Error JSON:\n%s", string(jsonData))

			expectedFields := []string{"name", "email", "age"}
			for _, field := range expectedFields {
				Expect(errorMap).To(HaveKey(field))
			}
		})
	})

	Context("Custom validator class", func() {
		It("should work with custom validator class", func() {
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
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("email"))
		})
	})

	Context("New validation rules", func() {
		It("should validate password strength", func() {
			validator := validation.NewValidatorWithDB(db)

			data := map[string]interface{}{
				"password": "Weak",
			}
			rules := map[string][]string{
				"password": {validation.Password()},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should validate UUID format", func() {
			validator := validation.NewValidatorWithDB(db)

			data := map[string]interface{}{
				"uuid": "invalid-uuid",
			}
			rules := map[string][]string{
				"uuid": {validation.UUID()},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})
})
