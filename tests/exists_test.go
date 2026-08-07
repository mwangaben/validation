package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/mwangaben/validation"
)

func TestChecker(t *testing.T) {
	RegisterFailHand
	RunSpecs(t, "Permission Checker Suite")
}

var _ = Describe("Exists Validation", func() {
	var validator *validation.Validator

	BeforeEach(func() {
		validator = validation.NewValidatorWithDB(db)
	})

	Context("Single field exists validation", func() {
		It("should validate existing name", func() {
			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject non-existing name", func() {
			data := map[string]interface{}{
				"name": "NonExistentUser",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("invalid"))
		})

		It("should validate existing email", func() {
			data := map[string]interface{}{
				"email": "john@example.com",
			}
			rules := map[string][]string{
				"email": {"exists:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject non-existing email", func() {
			data := map[string]interface{}{
				"email": "nonexistent@example.com",
			}
			rules := map[string][]string{
				"email": {"exists:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should validate existing product code", func() {
			data := map[string]interface{}{
				"code": "P001",
			}
			rules := map[string][]string{
				"code": {"exists:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject non-existing product code", func() {
			data := map[string]interface{}{
				"code": "P999",
			}
			rules := map[string][]string{
				"code": {"exists:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Exists with multiple conditions", func() {
		It("should validate multiple fields", func() {
			data := map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
			}
			rules := map[string][]string{
				"name":  {"exists:users,name"},
				"email": {"exists:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})
	})

	Context("Exists edge cases", func() {
		It("should reject empty string", func() {
			data := map[string]interface{}{
				"name": "",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should reject nil value", func() {
			data := map[string]interface{}{
				"name": nil,
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should validate numeric value as string", func() {
			data := map[string]interface{}{
				"age": "30",
			}
			rules := map[string][]string{
				"age": {"exists:users,age"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should reject non-existing numeric value", func() {
			data := map[string]interface{}{
				"age": 99,
			}
			rules := map[string][]string{
				"age": {"exists:users,age"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Exists with different tables", func() {
		It("should check in users table", func() {
			data := map[string]interface{}{
				"name": "Jane Smith",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should check in products table", func() {
			data := map[string]interface{}{
				"code": "P002",
			}
			rules := map[string][]string{
				"code": {"exists:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should fail cross-table check", func() {
			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"exists:products,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Exists combined with other rules", func() {
		It("should validate exists and email", func() {
			data := map[string]interface{}{
				"email": "john@example.com",
			}
			rules := map[string][]string{
				"email": {"required", "email", "exists:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should fail with invalid email", func() {
			data := map[string]interface{}{
				"email": "notanemail",
			}
			rules := map[string][]string{
				"email": {"required", "email", "exists:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("email"))
		})

		It("should fail with non-existing value", func() {
			data := map[string]interface{}{
				"name": "NonExistentUser",
			}
			rules := map[string][]string{
				"name": {"required", "string", "min:2", "exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})
	})

	Context("Case insensitive check", func() {
		It("should handle case insensitivity", func() {
			// In MySQL with utf8mb4_unicode_ci, this should be case-insensitive
			data := map[string]interface{}{
				"name": "john doe", // lowercase version of "John Doe"
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			// This may pass or fail depending on collation
			GinkgoWriter.Printf("Case-insensitive check result: %v", result)
			if !result {
				GinkgoWriter.Printf("Error: %s", validator.Error())
			}
		})
	})
})
