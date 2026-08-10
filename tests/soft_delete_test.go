package tests

import (
	"gorm.io/gorm"
	"testing"

	"github.com/mwangaben/factory/factory"
	"github.com/mwangaben/validation"
	"github.com/mwangaben/validation/examples/models"
	"github.com/mwangaben/validation/tests/helpers"
	. "github.com/onsi/gomega"
	_ "gorm.io/gorm"
)

var testDB *gorm.DB

// Helper to create a test context with Gomega registered
func runTest(t *testing.T, testFunc func()) {
	RegisterTestingT(t)
	testFunc()
}

func setupSoftDeleteTest() {
	var err error
	testDB, err = helpers.InitDB()
	Expect(err).NotTo(HaveOccurred())

	err = factory.NewDatabaseHelper(testDB).RefreshDatabase(&models.User{})
	Expect(err).NotTo(HaveOccurred())

	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}

	for i := range users {
		err := testDB.Create(&users[i]).Error
		Expect(err).ToNot(HaveOccurred())
	}

	err = testDB.Delete(&models.User{}, "name = ?", "Jane Smith").Error
	Expect(err).ToNot(HaveOccurred())
}

func TestSoftDeleteValidation(t *testing.T) {
	RegisterTestingT(t)
	setupSoftDeleteTest()

	t.Run("exists rule should NOT find soft-deleted users", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)

			data := map[string]interface{}{
				"name": "Jane Smith",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("invalid"))
		})
	})

	t.Run("exists rule should find non-deleted users", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)

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

	t.Run("exists_without_soft_delete rule should find soft-deleted users", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)

			data := map[string]interface{}{
				"name": "Jane Smith",
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})
	})

	t.Run("exists_without_soft_delete rule should find non-deleted users", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)

			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})
	})

	t.Run("numeric ID with soft delete", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)

			var user models.User
			testDB.Unscoped().Where("name = ?", "Jane Smith").First(&user)
			Expect(user.ID).ToNot(BeZero())

			// Test exists rule (should be false)
			data1 := map[string]interface{}{
				"id": user.ID,
			}
			rules1 := map[string][]string{
				"id": {"exists:users,id"},
			}
			result1 := validator.Validate(data1, rules1)
			Expect(result1).To(BeFalse())

			// Test exists_without_soft_delete rule (should be true)
			data2 := map[string]interface{}{
				"id": user.ID,
			}
			rules2 := map[string][]string{
				"id": {"exists_without_soft_delete:users,id"},
			}
			result2 := validator.Validate(data2, rules2)
			Expect(result2).To(BeTrue())
		})
	})
}
