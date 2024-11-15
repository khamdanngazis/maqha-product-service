// internal/handler/product_handler.go

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/product_service/internal/app/model"
	"maqhaa/product_service/internal/app/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// ProductHandler handles HTTP requests related to products.
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProductGroupsByCategoryHandler handles the GET request to fetch product groups by category.
func (h *ProductHandler) GetProductGroupsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	var appError service.AppError
	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	categories, appError := h.productService.GetProductGroupsByCategory(r.Context(), token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, categories)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var request *model.ProductCategoryRequest
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	appError = h.productService.AddProductCategoryService(r.Context(), request, token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) EditCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var request *model.ProductCategoryRequest
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["categoryID"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	request.ID = uint(categoryID)

	appError = h.productService.EditProductCategoryService(r.Context(), request, token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) DeactiveCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["categoryID"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	appError = h.productService.DeleteProductCategoryService(r.Context(), uint(categoryID), token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var request *model.ProductRequest
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	appError = h.productService.AddProductService(r.Context(), request, token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) EditProductHandler(w http.ResponseWriter, r *http.Request) {
	var request *model.ProductRequest
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Infof("Invalid request payload %s", err.Error())

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productID"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload productID")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	request.ID = uint(productID)

	appError = h.productService.EditProductService(r.Context(), request, token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *ProductHandler) DeactiveProductHandler(w http.ResponseWriter, r *http.Request) {
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productID"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	appError = h.productService.DeleteProductService(r.Context(), uint(productID), token)

	// Respond with the fetched categories
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}
