package tests

import (
	"fmt"
	"github.com/mwangaben/factory/factory"
	"github.com/mwangaben/validation"
	"github.com/mwangaben/validation/examples/models"
	"github.com/mwangaben/validation/tests/helpers"
	"github.com/mwangaben/validation/tests/helpers/factories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"testing"
)

// SoftDeleteUser extends User with DeletedAt
type SoftDeleteUser struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex:idx_users_name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex:idx_users_email"`
	Age       int            `gorm:"type:int"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SoftDeleteUser) TableName() string {
	return "users"
}

func TeardownSoftDeleteTestDB() {
	if testDB != nil {
		err := factory.NewDatabaseHelper(testDB).DropAllTables()
		if err != nil {
			fmt.Printf("Dropping of tables failed")
			return
		}
	}
	if testSQLDB != nil {
		testSQLDB.Close()
	}
}

func TestSoftDeleteSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	var err error
	testDB, err = helpers.InitDB()
	if err != nil {
		_ = fmt.Errorf("the DB connection failed %v", err)
	}
	helper := factory.NewDatabaseHelper(testDB)

	err = helper.RefreshDatabase(&models.User{}, &models.Product{})
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	factories.InitUserFactory(testDB)
	factories.InitProductFactory(testDB)

	// Insert test data
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}
	for _, user := range users {
		if err := testDB.Create(&user).Error; err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Soft delete one user
	if err := testDB.Delete(&SoftDeleteUser{}, "name = ?", "Jane Smith").Error; err != nil {
		t.Fatalf("Failed to soft delete user: %v", err)
	}

	RunSpecs(t, "Soft Delete Validation Suite")

	TeardownTestDB()
	//if sqlDB, _ := testDB.DB(); sqlDB != nil {
	//	sqlDB.Close()
	//}
}

var _ = Describe("Soft Delete Validation", func() {
	var validator *validation.Validator
	var db *gorm.DB

	BeforeEach(func() {
		db = testDB
		_ = factory.NewDatabaseHelper(db).RefreshDatabase(&models.User{}, &models.Product{})
		validator = validation.NewValidatorWithDB(db)
		Expect(validator).ToNot(BeNil())

	})

	Context("Exists rule with soft delete", func() {
		It("should find non-deleted users", func() {
			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should NOT find soft-deleted users", func() {
			data := map[string]interface{}{
				"name": "Jane Smith", // This user is soft-deleted
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("invalid"))
		})

		It("should find non-existing users as false", func() {
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

	Context("ExistsWithoutSoftDelete rule", func() {
		It("should find non-deleted users", func() {
			data := map[string]interface{}{
				"name": "John Doe",
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should find soft-deleted users (includes them)", func() {
			data := map[string]interface{}{
				"name": "Jane Smith", // This user is soft-deleted
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should NOT find non-existing users", func() {
			data := map[string]interface{}{
				"name": "NonExistentUser",
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(validator.Error()).To(ContainSubstring("invalid"))
		})
	})

	Context("Comparing both rules", func() {
		It("should have different results for soft-deleted users", func() {
			// Test exists rule (should be false)
			data1 := map[string]interface{}{
				"name": "Jane Smith",
			}
			rules1 := map[string][]string{
				"name": {"exists:users,name"},
			}
			result1 := validator.Validate(data1, rules1)
			Expect(result1).To(BeFalse())

			// Test exists_without_soft_delete rule (should be true)
			data2 := map[string]interface{}{
				"name": "Jane Smith",
			}
			rules2 := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result2 := validator.Validate(data2, rules2)
			Expect(result2).To(BeTrue())

			// Verify the difference in behavior
			Expect(result1).ToNot(Equal(result2))
		})
	})

	Context("Edge cases for soft delete", func() {
		It("should handle nil values for exists rule", func() {
			data := map[string]interface{}{
				"name": nil,
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should handle nil values for exists_without_soft_delete rule", func() {
			data := map[string]interface{}{
				"name": nil,
			}
			rules := map[string][]string{
				"name": {"exists_without_soft_delete:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should handle empty string for exists rule", func() {
			data := map[string]interface{}{
				"name": "",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
		})

		It("should handle numeric IDs with soft delete", func() {
			data := map[string]interface{}{
				"id": 1, // John Doe's ID (not deleted)
			}
			rules := map[string][]string{
				"id": {"exists:users,id"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeTrue())
		})

		It("should handle numeric IDs for soft-deleted users", func() {
			// Get Jane Smith's ID
			var user SoftDeleteUser
			db.Unscoped().Where("name = ?", "Jane Smith").First(&user)
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
})
