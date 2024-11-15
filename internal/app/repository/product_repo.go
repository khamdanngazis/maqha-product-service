// internal/repository/product_repository.go

package repository

import (
	"context"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/product_service/internal/app/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProductRepository handles database interactions related to products.
type ProductRepository interface {
	GetProductGroupsByCategory(ctx context.Context, token string) ([]entity.ProductCategory, error)
	AddProductCategory(ctx context.Context, category *entity.ProductCategory) error
	EditProductCategory(ctx context.Context, category *entity.ProductCategory) error
	GetProductByID(ctx context.Context, productID uint, token string) (*entity.Product, error)
	AddProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	EditProduct(ctx context.Context, product *entity.Product) error
	GetProductCategoryByID(ctx context.Context, productCategoryID uint) (*entity.ProductCategory, error)
	DeactivateProductCategory(ctx context.Context, ID uint) error
	DeactivateProduct(ctx context.Context, ID uint) error
}

// Implement the interface in the ProductRepository struct
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository instance.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetProductCategoryByID(ctx context.Context, productCategoryID uint) (*entity.ProductCategory, error) {
	var category entity.ProductCategory
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	if err := r.db.Where("id = ?", productCategoryID).First(&category).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductCategoryByID: %s", err.Error())
		return nil, err
	}
	return &category, nil
}

func (r *productRepository) AddProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	if err := r.db.Create(product).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error create category  %s", err.Error())
		return nil, err
	}
	return product, nil
}

func (r *productRepository) EditProduct(ctx context.Context, product *entity.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID uint, token string) (*entity.Product, error) {
	var product entity.Product
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	if err := r.db.Joins("JOIN product_category ON product_category.id = product.category_id").
		Joins("JOIN client ON client.id = product_category.client_id").
		Where("product.id = ? AND client.token = ?", productID, token).
		First(&product).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductGroupsByCategory  %s", err.Error())
		return nil, err
	}

	return &product, nil
}

// GetProductGroupsByCategoryAndClientID fetches product categories with associated products,
// grouped by category and filtered by client ID.
func (r *productRepository) GetProductGroupsByCategory(ctx context.Context, token string) ([]entity.ProductCategory, error) {
	var categories []entity.ProductCategory
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)

	if err := r.db.
		Preload("Products").
		Joins("LEFT JOIN product ON product_category.id = product.category_id").
		Joins("LEFT JOIN client ON client.id = product_category.client_id").
		Where("client.token = ?", token).Order("product_category.id asc").
		Find(&categories).
		Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductGroupsByCategory  %s", err.Error())
		return nil, err
	}

	return categories, nil
}

func (r *productRepository) AddProductCategory(ctx context.Context, category *entity.ProductCategory) error {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	result := r.db.Create(category)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error create category  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

// EditProductCategory edits an existing ProductCategory in the database.
func (r *productRepository) EditProductCategory(ctx context.Context, category *entity.ProductCategory) error {
	result := r.db.Model(category).Updates(category)
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error create category  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

func (r *productRepository) DeactivateProductCategory(ctx context.Context, ID uint) error {

	logID, _ := ctx.Value(middleware.RequestIDKey).(string)

	updates := map[string]interface{}{
		"IsActive": false,
	}

	// Perform the update operation
	result := r.db.Model(&entity.ProductCategory{}).Where("id = ?", ID).Updates(updates)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error UpdateUser  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

func (r *productRepository) DeactivateProduct(ctx context.Context, ID uint) error {

	logID, _ := ctx.Value(middleware.RequestIDKey).(string)

	updates := map[string]interface{}{
		"IsActive": false,
	}

	// Perform the update operation
	result := r.db.Model(&entity.Product{}).Where("id = ?", ID).Updates(updates)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error UpdateUser  %s", result.Error.Error())
		return result.Error
	}
	return nil
}
