package tests

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mwangaben/validation/tests/helpers"
	"github.com/mwangaben/validation/tests/helpers/factories"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mwangaben/factory/factory"
	"github.com/mwangaben/validation"
	"github.com/mwangaben/validation/examples/models"
	"gorm.io/gorm"
)

var (
	testDB    *gorm.DB
	testSQLDB *sql.DB
)

func SetupTestDB() error {
	var err error
	testDB, err = helpers.InitDB()
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	testSQLDB, err = testDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}

	// Refresh database with all models
	helper := factory.NewDatabaseHelper(testDB)
	if err := helper.RefreshDatabase(&models.User{}, &models.Product{}); err != nil {
		return fmt.Errorf("failed to refresh: %v", err)
	}

	// Initialize factories
	factories.InitUserFactory(testDB)
	factories.InitProductFactory(testDB)

	// Seed with random data
	_, err = factories.UserFactory().ClearOverrides().Count(20).CreateMany()
	if err != nil {
		return fmt.Errorf("failed to create users: %v", err)
	}

	_, err = factories.ProductFactory().ClearOverrides().Count(10).CreateMany()
	if err != nil {
		return fmt.Errorf("failed to create products: %v", err)
	}

	fmt.Println("✅ Test data seeded successfully")
	return nil
}

func TeardownTestDB() {
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

func TestExistsSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	if err := SetupTestDB(); err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}
	ginkgo.RunSpecs(t, "Exists Validation Suite")
	TeardownTestDB()
}

// ============ TESTS ============

var _ = ginkgo.Describe("Exists Validation", func() {
	var validator *validation.Validator

	ginkgo.BeforeEach(func() {
		validator = validation.NewValidatorWithDB(testDB)
		gomega.Expect(validator).ToNot(gomega.BeNil())
	})

	ginkgo.Context("Single field exists validation", func() {
		ginkgo.It("should validate existing name", func() {
			var user models.User
			err := testDB.First(&user).Error
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			data := map[string]interface{}{
				"name": user.Name,
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			fmt.Printf("The User is %v", user)
			result := validator.Validate(data, rules)
			gomega.Expect(result).To(gomega.BeTrue())
		})

		ginkgo.It("should reject non-existing name", func() {
			data := map[string]interface{}{
				"name": "NonExistentUser_XYZ123",
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			gomega.Expect(result).To(gomega.BeFalse())
		})
	})

	ginkgo.Context("Exists with custom data using factories", func() {
		ginkgo.It("should validate user created with custom data", func() {
			uniqueName := fmt.Sprintf("Custom_User_%d", gofakeit.Number(100000, 999999))
			uniqueEmail := fmt.Sprintf("custom_%d@example.com", gofakeit.Number(100000, 999999))

			user, err := factories.UserFactory().
				WithOverrides(map[string]interface{}{
					"name":  uniqueName,
					"email": uniqueEmail,
					"age":   25,
				}).
				Create()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(user.Name).To(gomega.Equal(uniqueName))

			data := map[string]interface{}{
				"name": uniqueName,
			}
			rules := map[string][]string{
				"name": {"exists:users,name"},
			}
			result := validator.Validate(data, rules)
			gomega.Expect(result).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("Product validation", func() {
		ginkgo.It("should validate existing product code", func() {
			var product models.Product
			err := testDB.First(&product).Error
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			data := map[string]interface{}{
				"code": product.Code,
			}
			rules := map[string][]string{
				"code": {"exists:products,code"},
			}
			result := validator.Validate(data, rules)
			gomega.Expect(result).To(gomega.BeTrue())
		})

		ginkgo.It("should validate product with custom code", func() {
			uniqueCode := fmt.Sprintf("P%d", gofakeit.Number(100000, 999999))

			_, err := factories.ProductFactory().
				WithOverrides(map[string]interface{}{
					"code":  uniqueCode,
					"name":  "Custom Product",
					"price": 199.99,
				}).
				Create()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			data := map[string]interface{}{
				"code": uniqueCode,
			}
			rules := map[string][]string{
				"code": {"exists:products,code"},
			}
			result := validator.Validate(data, rules)
			gomega.Expect(result).To(gomega.BeTrue())
		})
	})
})
