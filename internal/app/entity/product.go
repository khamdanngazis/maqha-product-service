package entity

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"categoryId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Set the table name explicitly for GORM
func (Product) TableName() string {
	return "product"
}
