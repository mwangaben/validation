package tests

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/mwangaben/factory/factory"
	"github.com/mwangaben/validation"
	"github.com/mwangaben/validation/examples/models"
	"github.com/mwangaben/validation/tests/helpers"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

// ============ GLOBAL SETUP ============

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func init() {
	RegisterFailHandler(func(message string, callerSkip ...int) {
		panic(message)
	})
}

// ============ SETUP ============

func setupExistTest() {
	var err error
	db, err = helpers.InitDBPost() // Fixed: assign to global db (no colon)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	sqlDB, err = db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get SQL DB: %v", err))
	}

	err = factory.NewDatabaseHelper(db).RefreshDatabase(&models.User{}, &models.Product{})
	if err != nil {
		panic(fmt.Sprintf("Failed to refresh database: %v", err))
	}

	// Insert test data
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			panic(fmt.Sprintf("Failed to insert user test data: %v", err))
		}
	}

	products := []models.Product{
		{Code: "P001", Name: "Laptop", Price: 999.99},
		{Code: "P002", Name: "Mouse", Price: 29.99},
	}
	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			panic(fmt.Sprintf("Failed to insert product test data: %v", err))
		}
	}

	fmt.Println("✅ Test data inserted successfully")
}

func teardownExistTest() {
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// ============ TESTS ============

func TestExistsValidation(t *testing.T) {
	setupExistTest()
	defer teardownExistTest()

	runTest(t, func() {
		validator := validation.NewValidatorWithDB(db)
		Expect(validator).ToNot(BeNil())

		t.Run("should validate existing name", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"name": "John Doe",
				}
				rules := map[string][]string{
					"name": {"exists:users,name"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject non-existing name", func(t *testing.T) {
			runTest(t, func() {
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
		})

		t.Run("should validate existing email", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"email": "john@example.com",
				}
				rules := map[string][]string{
					"email": {"exists:users,email"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject non-existing email", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"email": "nonexistent@example.com",
				}
				rules := map[string][]string{
					"email": {"exists:users,email"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate existing product code", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"code": "P001",
				}
				rules := map[string][]string{
					"code": {"exists:products,code"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject non-existing product code", func(t *testing.T) {
			runTest(t, func() {
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

		t.Run("should validate multiple fields", func(t *testing.T) {
			runTest(t, func() {
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

		t.Run("should reject empty string", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"name": "",
				}
				rules := map[string][]string{
					"name": {"exists:users,name"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should reject nil value", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"name": nil,
				}
				rules := map[string][]string{
					"name": {"exists:users,name"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeFalse())
			})
		})

		t.Run("should validate numeric value as string", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"age": "30",
				}
				rules := map[string][]string{
					"age": {"exists:users,age"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should reject non-existing numeric value", func(t *testing.T) {
			runTest(t, func() {
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

		t.Run("should check in users table", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"name": "Jane Smith",
				}
				rules := map[string][]string{
					"name": {"exists:users,name"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should check in products table", func(t *testing.T) {
			runTest(t, func() {
				data := map[string]interface{}{
					"code": "P002",
				}
				rules := map[string][]string{
					"code": {"exists:products,code"},
				}
				result := validator.Validate(data, rules)
				Expect(result).To(BeTrue())
			})
		})

		t.Run("should fail cross-table check", func(t *testing.T) {
			runTest(t, func() {
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
	})
}
