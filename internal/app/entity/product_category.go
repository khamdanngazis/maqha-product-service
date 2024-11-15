package entity

import (
	"time"
)

type ProductCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClientID  uint      `json:"clientId"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"products"`
}

// Set the table name explicitly for GORM
func (ProductCategory) TableName() string {
	return "product_category"
}
