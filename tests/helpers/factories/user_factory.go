package factories

import (
	"github.com/mwangaben/factory/factory"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mwangaben/validation/examples/models"
	"gorm.io/gorm"
)

// productFactory is the internal user factories instance
var userFactpry *factory.ModelFactory[models.User]

// InitUserFactory initializes the user factories
func InitUserFactory(db *gorm.DB) {
	userFactpry = factory.NewModelFactory(db, func() models.User {
		return models.User{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Age:   gofakeit.Age(),
		}
	})
}

// UserFactory returns the user factories instance
func UserFactory() *factory.ModelFactory[models.User] {
	if userFactpry == nil {
		panic("User factories not initialized. Call factories.Init() first.")
	}
	return userFactpry
}
