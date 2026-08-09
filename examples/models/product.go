package models

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"type:varchar(50);uniqueIndex:idx_products_code"`
	Name      string `gorm:"type:varchar(100)"`
	Price     float64
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (Product) TableName() string {
	return "products"
}
