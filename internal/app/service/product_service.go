// internal/service/product_service.go

package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"maqhaa/library/helper"
	exRepo "maqhaa/product_service/external/repository"
	"maqhaa/product_service/internal/app/entity"
	"maqhaa/product_service/internal/app/model"
	"maqhaa/product_service/internal/app/repository"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ProductService handles business logic related to products.
type ProductService interface {
	GetProductGroupsByCategory(ctx context.Context, token string) ([]entity.ProductCategory, AppError)
	GetProductByID(ctx context.Context, ID uint, token string) (*entity.Product, AppError)
	AddProductCategoryService(ctx context.Context, request *model.ProductCategoryRequest, token string) AppError
	EditProductCategoryService(ctx context.Context, request *model.ProductCategoryRequest, token string) AppError
	DeleteProductCategoryService(ctx context.Context, ID uint, token string) AppError
	AddProductService(ctx context.Context, request *model.ProductRequest, token string) AppError
	EditProductService(ctx context.Context, request *model.ProductRequest, token string) AppError
	DeleteProductService(ctx context.Context, ID uint, token string) AppError
}

// productServiceImpl implements the ProductService interface
type productServiceImpl struct {
	productRepository repository.ProductRepository
	userRepository    exRepo.UserRepository
	imageRepository   repository.ImagesRepository
}

// NewProductService creates a new ProductService instance.
func NewProductService(productRepository repository.ProductRepository, userRepository exRepo.UserRepository, imageRepository repository.ImagesRepository) ProductService {
	return &productServiceImpl{
		productRepository: productRepository,
		userRepository:    userRepository,
		imageRepository:   imageRepository,
	}
}

// GetProductGroupsByCategory fetches product groups (categories with associated products),
// grouped by category and filtered by token.
func (s *productServiceImpl) GetProductGroupsByCategory(ctx context.Context, token string) ([]entity.ProductCategory, AppError) {
	if token == "" {
		return nil, *NewInvalidTokenError()
	}

	categories, err := s.productRepository.GetProductGroupsByCategory(ctx, token)
	if err != nil && err.Error() != "record not found" {
		return nil, *NewQueryDBError()
	}
	encountered := map[uint]bool{}
	result := []entity.ProductCategory{}

	for _, category := range categories {
		if !encountered[category.ID] {
			encountered[category.ID] = true
			result = append(result, category)
		}
	}

	if len(result) == 0 {
		return nil, *NewInvalidTokenError()
	}
	return result, *NewSuccessError()
}

func (s *productServiceImpl) GetProductByID(ctx context.Context, ID uint, token string) (*entity.Product, AppError) {
	product, err := s.productRepository.GetProductByID(ctx, ID, token)

	if err != nil {

		if err.Error() != "record not found" {
			return nil, *NewQueryDBError()
		}
	}

	if product == nil {
		return nil, *NewProductNotFoundError()
	}

	return product, *NewSuccessError()
}

func (s *productServiceImpl) AddProductCategoryService(ctx context.Context, request *model.ProductCategoryRequest, token string) AppError {

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}
	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	category := &entity.ProductCategory{
		ClientID: uint(user.ClientId),
		Name:     request.Category,
	}
	err = s.productRepository.AddProductCategory(ctx, category)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

// EditProductCategoryService edits an existing ProductCategory.
func (s *productServiceImpl) EditProductCategoryService(ctx context.Context, request *model.ProductCategoryRequest, token string) AppError {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}
	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	category, err := s.productRepository.GetProductCategoryByID(ctx, request.ID)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if category.ClientID != uint(user.ClientId) {
		return *NewInvalidTokenError()
	}

	category.Name = request.Category

	err = s.productRepository.EditProductCategory(ctx, category)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

// DeleteProductCategoryService deletes a ProductCategory.
func (s *productServiceImpl) DeleteProductCategoryService(ctx context.Context, ID uint, token string) AppError {

	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	category, err := s.productRepository.GetProductCategoryByID(ctx, ID)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if category.ClientID != uint(user.ClientId) {
		return *NewInvalidTokenError()
	}

	err = s.productRepository.DeactivateProductCategory(ctx, ID)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (s *productServiceImpl) AddProductService(ctx context.Context, request *model.ProductRequest, token string) AppError {

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}
	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	category, err := s.productRepository.GetProductCategoryByID(ctx, request.CategoryID)

	if err != nil {
		return *NewQueryDBError()
	}

	if category.ClientID != uint(user.ClientId) {
		return *NewInvalidTokenError()
	}

	productImage, err := SaveImage(request.Image, s.imageRepository)

	if err != nil {
		return *NewInvalidRequestError(err.Error())
	}

	product := &entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Image:       productImage,
		CategoryID:  request.CategoryID,
	}

	_, err = s.productRepository.AddProduct(ctx, product)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (s *productServiceImpl) EditProductService(ctx context.Context, request *model.ProductRequest, token string) AppError {

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}
	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	category, err := s.productRepository.GetProductCategoryByID(ctx, request.CategoryID)

	if err != nil {
		return *NewQueryDBError()
	}

	if category.ClientID != uint(user.ClientId) {
		return *NewInvalidTokenError()
	}

	product, err := s.productRepository.GetProductByID(ctx, request.ID, token)
	if err != nil {
		return *NewInvalidTokenError()
	}

	productImage, err := SaveImage(request.Image, s.imageRepository)
	if err != nil {
		return *NewInvalidRequestError(err.Error())
	}

	err = s.imageRepository.RemoveImage(product.Image)
	if err != nil {
		fmt.Println(err.Error())
	}

	updateProduct := &entity.Product{
		ID:          product.ID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Image:       productImage,
		CategoryID:  product.CategoryID,
	}

	err = s.productRepository.EditProduct(ctx, updateProduct)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (s *productServiceImpl) DeleteProductService(ctx context.Context, ID uint, token string) AppError {

	user, err := s.userRepository.GetUser(ctx, token)

	if err != nil {
		return *NewInvalidTokenError()
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewInvalidTokenError()
	}

	product, err := s.productRepository.GetProductByID(ctx, ID, token)
	if err != nil {
		return *NewInvalidTokenError()
	}

	category, err := s.productRepository.GetProductCategoryByID(ctx, product.CategoryID)

	if err != nil {
		return *NewQueryDBError()
	}

	if category.ClientID != uint(user.ClientId) {
		return *NewInvalidTokenError()
	}

	err = s.productRepository.DeactivateProduct(ctx, ID)

	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func SaveImage(base64Image string, imagesRepo repository.ImagesRepository) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(imageData)
	mimeTypePrefix := strings.Split(mimeType, "/")[1]

	id := "PR"
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	imageName := fmt.Sprintf("%s-%d.%s", id, ms, mimeTypePrefix)

	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}
	compressedImage, err := helper.CompressImage(img, format)
	if err != nil {
		return "", err
	}

	err = imagesRepo.SaveImage(compressedImage, imageName)
	if err != nil {
		return "", err
	}

	return imageName, nil
}
