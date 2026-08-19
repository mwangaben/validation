package tests

import (
	"database/sql"
	"fmt"
	"github.com/mwangaben/factory/factory"
	"github.com/mwangaben/validation"
	"github.com/mwangaben/validation/examples/models"
	"github.com/mwangaben/validation/tests/helpers"
	"github.com/mwangaben/validation/tests/helpers/factories"
	. "github.com/onsi/gomega"
	"testing"
)

var testSQLDB *sql.DB

func SetUpTestDB() error {
	var err error
	testDB, err = helpers.InitDBPost()
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

	// Initialize factories
	factories.InitUserFactory(testDB)
	factories.InitProductFactory(testDB)
	return nil
}

func TeardownTestDB() {
	if testSQLDB != nil {
		testSQLDB.Close()
	}
}

func TestUniqueSuite(t *testing.T) {
	RegisterTestingT(t)

	if err := SetUpTestDB(); err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}
	defer TeardownTestDB()

	t.Run("Single field unique validation", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)
			Expect(validator).ToNot(BeNil())

			t.Run("should validate unique name", func(t *testing.T) {
				runTest(t, func() {
					data := map[string]interface{}{
						"name": "UniqueName",
					}
					rules := map[string][]string{
						"name": {"unique:users,name"},
					}
					result := validator.Validate(data, rules)
					Expect(result).To(BeTrue())
				})
			})

			t.Run("should reject duplicate name", func(t *testing.T) {
				runTest(t, func() {
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
			})

			t.Run("should validate unique email", func(t *testing.T) {
				runTest(t, func() {
					data := map[string]interface{}{
						"email": "newuser@example.com",
					}
					rules := map[string][]string{
						"email": {"unique:users,email"},
					}
					result := validator.Validate(data, rules)
					Expect(result).To(BeTrue())
				})
			})

			t.Run("should reject duplicate email", func(t *testing.T) {
				runTest(t, func() {
					user, err := factories.UserFactory().WithOverrides(map[string]interface{}{
						"name":  "Benedict",
						"email": "john@example.com",
					}).Create()
					Expect(err).NotTo(HaveOccurred())
					Expect(user.Name).To(Equal("Benedict"))

					data := map[string]interface{}{
						"email": "john@example.com",
					}
					rules := map[string][]string{
						"email": {"unique:users,email"},
					}
					result := validator.Validate(data, rules)
					Expect(result).To(BeFalse())
				})
			})

			t.Run("should validate unique product code", func(t *testing.T) {
				runTest(t, func() {
					data := map[string]interface{}{
						"code": "P003",
					}
					rules := map[string][]string{
						"code": {"unique:products,code"},
					}
					result := validator.Validate(data, rules)
					Expect(result).To(BeTrue())
				})
			})

			t.Run("should reject duplicate product code", func(t *testing.T) {
				runTest(t, func() {
					product, err := factories.ProductFactory().WithOverrides(map[string]interface{}{
						"code": "P001",
					}).Create()
					Expect(err).NotTo(HaveOccurred())
					Expect(product.Code).To(Equal("P001"))

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
		})
	})

	t.Run("Unique with except ID", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)
			Expect(validator).ToNot(BeNil())

			t.Run("should pass when excepting own ID", func(t *testing.T) {
				runTest(t, func() {
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
			})

			t.Run("should fail when excepting different ID", func(t *testing.T) {
				runTest(t, func() {
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
		})
	})

	t.Run("Structured GraphQL errors", func(t *testing.T) {
		runTest(t, func() {
			validator := validation.NewValidatorWithDB(testDB)
			Expect(validator).ToNot(BeNil())

			t.Run("should return structured GraphQL errors", func(t *testing.T) {
				runTest(t, func() {
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
	})
}
