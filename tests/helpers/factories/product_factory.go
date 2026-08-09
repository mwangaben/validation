package factories

import (
	"fmt"
	"github.com/mwangaben/factory/factory"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mwangaben/validation/examples/models"
	"gorm.io/gorm"
)

// productFactory is the internal product factories instance
var productFactory *factory.ModelFactory[models.Product]

// InitProductFactory initializes the product factories
func InitProductFactory(db *gorm.DB) {
	productFactory = factory.NewModelFactory(db, func() models.Product {
		return models.Product{
			Code:  fmt.Sprintf("P%d", gofakeit.Number(100000, 999999)),
			Name:  fmt.Sprintf("%s %s", gofakeit.ProductName(), gofakeit.UUID()[0:4]),
			Price: gofakeit.Price(10, 1000),
		}
	})
}

// ProductFactory returns the product factories instance
func ProductFactory() *factory.ModelFactory[models.Product] {
	if productFactory == nil {
		panic("Product factories not initialized. Call factories.Init() first.")
	}
	return productFactory
}
