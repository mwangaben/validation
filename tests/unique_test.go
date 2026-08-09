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
	"testing"
)

func SetUpTestDB() error {
	var err error
	testDB, err = helpers.InitDB()
	if err != nil {
		return fmt.Errorf("the DB connection failed %v", err)
	}

	testSQLDB, err = testDB.DB()
	if err != nil {
		return fmt.Errorf("failed to load testDB.DB: %v", err)
	}

	helper := factory.NewDatabaseHelper(testDB)
	if err := helper.RefreshDatabase(&models.User{}, &models.Product{}); err != nil {
		return fmt.Errorf("the database refresh failed %v", err)
	}

	//initialize factories.
	factories.InitUserFactory(testDB)
	factories.InitProductFactory(testDB)
	return nil
}

func TestUniqueSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	if err := SetUpTestDB(); err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}
	RunSpecs(t, "Exists Validation Suite")
	TeardownTestDB()
}

var _ = Describe("Unique Validation", func() {
	var validator *validation.Validator

	BeforeEach(func() {
		validator = validation.NewValidatorWithDB(testDB)
		Expect(validator).ToNot(BeNil())
		//helper := factory.NewDatabaseHelper(testDB)
		//if err := helper.RefreshDatabase(&models.User{}); err != nil {
		//	_ = fmt.Errorf("the database refresh failed %v", err)
		//}
		//_, err := factories.UserFactory().WithOverrides(map[string]interface{}{
		//	"name": "John Doe",
		//}).Create()
		//if err != nil {
		//	fmt.Printf("Failed to run factories %v", err)
		//}
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
			user, err := factories.UserFactory().WithOverrides(data).Create()
			Expect(err).ToNot(HaveOccurred())
			Expect(user.Name).To(Equal(data["name"]))

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
			user, err := factories.UserFactory().WithOverrides(map[string]interface{}{
				"name":  "Benedict",
				"email": "john@example.com",
			}).Create()

			data := map[string]interface{}{
				"email": "john@example.com",
			}
			rules := map[string][]string{
				"email": {"unique:users,email"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(user.Name).To(Equal("Benedict"))
			Expect(err).NotTo(HaveOccurred())
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
			product, err := factories.ProductFactory().WithOverrides(map[string]interface{}{
				"code": "P001",
			}).Create()
			data := map[string]interface{}{
				"code": "P001",
			}
			rules := map[string][]string{
				"code": {"unique:products,code"},
			}
			result := validator.Validate(data, rules)
			Expect(result).To(BeFalse())
			Expect(product.Code).To(Equal("P001"))
			Expect(err).NotTo(HaveOccurred())
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

	Context("Structured GraphQL errors", func() {
		It("should return structured GraphQL errors", func() {
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
		})
	})
})
