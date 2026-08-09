package models

import "time"

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	ProductID uint `gorm:"index"`
	Quantity  int
	Total     float64
	CreatedAt time.Time
}

func (Order) TableName() string {
	return "orders"
}
