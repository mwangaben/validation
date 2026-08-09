package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);uniqueIndex:idx_users_name"`
	Email     string `gorm:"type:varchar(100);uniqueIndex:idx_users_email"`
	Age       int
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

func (User) TableName() string {
	return "users"
}
