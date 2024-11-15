package model

import "time"

type GetProductRequest struct {
	ID uint `json:"product_id"`
}

type GetProductResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		ID          uint      `json:"id"`
		CategoryID  uint      `json:"categoryId"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		Price       float64   `json:"price"`
		IsActive    bool      `json:"isActive"`
		CreatedAt   time.Time `json:"createdAt"`
	} `json:"data,omitempty"`
}

type ProductCategoryRequest struct {
	ID       uint
	Category string `json:"category" validate:"required"`
}

type ProductRequest struct {
	ID          uint
	CategoryID  uint    `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}
