// product_handler_test.go

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/product_service/internal/app/entity"
	"maqhaa/product_service/internal/app/model"
	"maqhaa/product_service/internal/app/service"

	pb "maqhaa/product_service/internal/interface/grpc/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	exModel "maqhaa/product_service/external/model"
)

func TestGetProductGroupsByCategoryHandler_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	categories := SampleCategories(client.ID)
	for _, category := range categories {
		db.Create(category)
	}

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Mock HTTP request
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the client token in the request header
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(productHandler.GetProductGroupsByCategoryHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response body
	var response model.HTTPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("canot Unmarshal response err : %v", err)
	}
	var dataCategories []entity.ProductCategory
	switch data := response.Data.(type) {
	case nil:
		t.Errorf("Data is nil")
	case []interface{}:
		for _, item := range data {
			categoryData, err := json.Marshal(item)
			if err != nil {
				t.Errorf("Error marshaling category data: %v", err)
				continue
			}

			var category entity.ProductCategory
			if err := json.Unmarshal(categoryData, &category); err != nil {
				t.Errorf("Error unmarshaling category data: %v", err)
				continue
			}

			dataCategories = append(dataCategories, category)
		}

	default:
		t.Errorf("Data has unexpected type: %T\n", response.Data)
	}

	if (len(dataCategories)) != len(categories) {
		t.Errorf("Unexpected number of categories. Got %d, want %d", len(dataCategories), len(categories))
	}

	// Perform assertions based on the expected data
	// You might want to compare responseCategories with the expected categories

	// For example, you can check if the number of categories returned is as expected
	/*
		if len(responseCategories) != len(categories) {
			t.Errorf("Unexpected number of categories. Got %d, want %d", len(responseCategories), len(categories))
		}*/

	// Add more assertions based on your application's logic

}

func TestGetProductGroupsByCategoryHandler_TokenIsNull(t *testing.T) {
	// Mock HTTP request with null token
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.GetProductGroupsByCategoryHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")

	// Check the response status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for null token: got %v want %v", status, http.StatusBadRequest)
	}

	// Parse the response body
	var response model.HTTPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Cannot unmarshal response: %v", err)
	}

	// Perform assertions based on the expected error response
	// For example, you can check if the error message is as expected
	if response.Message != "Invalid Token" {
		t.Errorf("Unexpected error message for null token: got %v want %v", response.Message, "Invalid Token")
	}
}

func TestGetProductGroupsByCategoryHandler_InvalidToken(t *testing.T) {
	// Mock HTTP request with an invalid token
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Set an invalid token in the request header
	req.Header.Set("Token", "invalid_token")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(productHandler.GetProductGroupsByCategoryHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid token: got %v want %v", status, http.StatusBadRequest)
	}
	// Parse the response body
	var response model.HTTPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Cannot unmarshal response: %v", err)
	}

	// Perform assertions based on the expected error response
	// For example, you can check if the error message indicates unauthorized access
	if response.Message != "Invalid Token" {
		t.Errorf("Unexpected error message for invalid token: got %v want %v", response.Message, "Invalid Token")
	}
}

func TestGetProductByIDGRPCHandler_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	categories := SampleCategories(client.ID)
	for _, category := range categories {
		db.Create(category)
	}

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Set up a gRPC connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error creating gRPC client connection: %v", err)
	}
	defer conn.Close()

	clientServer := pb.NewProductClient(conn)

	// Prepare a request
	req := &pb.GetProductRequest{
		ProductId: uint32(categories[0].Products[0].ID),
		Token:     client.Token, // Replace with a valid product ID for your test data
	}

	// Call the gRPC method
	resp, err := clientServer.GetProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Error calling GetProduct gRPC method: %v", err)
	}

	// Assertions
	assert.NotNil(t, resp)
	assert.Equal(t, int32(service.SuccessError), resp.Code) // Assuming 0 is the success code
	assert.Equal(t, service.SuccessMessage, resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, categories[0].Products[0].Name, resp.Data.Name)
	assert.Equal(t, float32(categories[0].Products[0].Price), resp.Data.Price)
	assert.Equal(t, categories[0].Products[0].Description, resp.Data.Description)
	// Add more assertions based on your response structure

	// Cleanup if needed
	// ...

}

func TestGetProductByIDGRPCHandler_InvalidClient(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	categories := SampleCategories(client.ID)
	for _, category := range categories {
		db.Create(category)
	}

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Set up a gRPC connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error creating gRPC client connection: %v", err)
	}
	defer conn.Close()

	clientServer := pb.NewProductClient(conn)

	// Prepare a request
	req := &pb.GetProductRequest{
		ProductId: uint32(categories[0].Products[0].ID),
		Token:     "invalid_token", // Replace with a valid product ID for your test data
	}

	// Call the gRPC method
	resp, err := clientServer.GetProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Error calling GetProduct gRPC method: %v", err)
	}

	// Assertions
	assert.NotNil(t, resp)
	assert.Equal(t, int32(service.ProductNotFound), resp.Code) // Assuming 0 is the success code
	assert.Equal(t, service.ProductNotFoundMessage, resp.Message)
	assert.Nil(t, resp.Data)
	// Add more assertions based on your response structure

	// Cleanup if needed
	// ...

}

func TestGetProductByIDGRPCHandler_ProductNotFound(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	categories := SampleCategories(client.ID)
	for _, category := range categories {
		db.Create(category)
	}

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Set up a gRPC connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error creating gRPC client connection: %v", err)
	}
	defer conn.Close()

	clientServer := pb.NewProductClient(conn)

	// Prepare a request
	req := &pb.GetProductRequest{
		ProductId: 0,
		Token:     client.Token, // Replace with a valid product ID for your test data
	}

	// Call the gRPC method
	resp, err := clientServer.GetProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Error calling GetProduct gRPC method: %v", err)
	}

	// Assertions
	assert.NotNil(t, resp)
	assert.Equal(t, int32(service.ProductNotFound), resp.Code) // Assuming 0 is the success code
	assert.Equal(t, service.ProductNotFoundMessage, resp.Message)
	assert.Nil(t, resp.Data)
	// Add more assertions based on your response structure

	// Cleanup if needed
	// ...

}

func TestAddProductCategory_Success(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/category", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.AddCategoryHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var category entity.ProductCategory
	result := db.First(&category)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	assert.Equal(t, category.Name, "New Category")
}

func TestAddProductCategory_InvalidParam(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/category", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.AddCategoryHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidRequestError, response.Code)
}

func TestEditProductCategory_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.EditCategoryHandler).Methods("PUT")

	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/category/"+strconv.Itoa(int(categories[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var category entity.ProductCategory
	result := db.First(&category)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	assert.Equal(t, category.Name, "New Category")
}

func TestEditProductCategory_InvalidToken(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.EditCategoryHandler).Methods("PUT")

	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/category/"+strconv.Itoa(int(categories[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", "")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestEditProductCategory_CilentIDNotMatch(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	client2 := SampleClient2()
	db.Create(client2)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client2.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.EditCategoryHandler).Methods("PUT")

	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/category/"+strconv.Itoa(int(categories[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestDeactiveProductCategory_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.DeactiveCategoryHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/category/"+strconv.Itoa(int(categories[0].ID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var category entity.ProductCategory
	result := db.First(&category)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	fmt.Println(category)

	assert.Equal(t, false, category.IsActive)
}

func TestDeactiveProductCategory_CilentIDNotMatch(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	client2 := SampleClient2()
	db.Create(client2)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client2.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.EditCategoryHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/category/"+strconv.Itoa(int(categories[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestDeactiveProductCategory_InvalidToken(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	client2 := SampleClient2()
	db.Create(client2)
	token := "xxxxxaaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client2.ID)
	categories[0].ID = 1
	db.Create(categories[0])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductCategoryRequest{
		Category: "New Category",
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/category/{categoryID}", productHandler.EditCategoryHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/category/"+strconv.Itoa(int(categories[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", "token")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestAddProduct_Success(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxabaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[3].ID = 1
	db.Create(categories[3])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductRequest{
		Name:        "Mochacino",
		CategoryID:  categories[3].ID,
		Description: "Mochacino",
		Price:       25000,
		Image:       SampleImagePNG(),
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(string(requestJSON[:]))

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.AddProductHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var product entity.Product
	result := db.First(&product)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	assert.Equal(t, product.Name, request.Name)
	assert.Equal(t, request.Description, product.Description)
	assert.Equal(t, request.Price, product.Price)
	assert.Equal(t, request.Description, product.Description)
	assert.Equal(t, request.CategoryID, product.CategoryID)
}

func TestAddProduct_InvalidImage(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxabaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[3].ID = 1
	db.Create(categories[3])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductRequest{
		Name:        "Mochacino",
		CategoryID:  categories[3].ID,
		Description: "Mochacino",
		Price:       25000,
		Image:       SampleImageGif(),
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(string(requestJSON[:]))

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.AddProductHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.InvalidRequestError, response.Code)

}

func TestAddProduct_InvalidClientID(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	token := "xxxxxabaaaa"
	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	client2 := SampleClient2()
	db.Create(client2)
	categories := SampleCategories(client2.ID)
	categories[3].ID = 1
	db.Create(categories[3])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductRequest{
		Name:        "Mochacino",
		CategoryID:  categories[3].ID,
		Description: "Mochacino",
		Price:       25000,
		Image:       SampleImageGif(),
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(string(requestJSON[:]))

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(productHandler.AddProductHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestEditProduct_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	token := "xxxxxaaaaa"
	client.Token = token
	db.Create(client)

	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[2].ID = 1
	db.Create(categories[2])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	// Create a login request
	request := model.ProductRequest{
		Name:        "New Product",
		Description: "new description",
		CategoryID:  categories[2].ID,
		Price:       100000,
		Image:       SampleImagePNG(),
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(string(requestJSON[:]))

	router := mux.NewRouter()
	router.HandleFunc("/product/{productID}", productHandler.EditProductHandler).Methods("PUT")

	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/product/"+strconv.Itoa(int(categories[2].Products[0].ID)), bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var product entity.Product
	result := db.First(&product)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	assert.Equal(t, request.Name, product.Name)

}

func TestDeactiveProduct_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	token := "xxxxxaaaaa"
	client.Token = token
	db.Create(client)

	userRepo.SetUserResponse(token, &exModel.UserData{Id: 1, ClientId: uint32(client.ID), IsAdmin: true, IsLogin: true})

	categories := SampleCategories(client.ID)
	categories[2].ID = 1
	db.Create(categories[2])

	// Clean up the testing environment
	tables := []string{"product", "product_category", "client"}
	defer clearDB(tables)

	router := mux.NewRouter()
	router.HandleFunc("/product/{productID}", productHandler.DeactiveProductHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/product/"+strconv.Itoa(int(categories[2].Products[0].ID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	// Gorm query to get the first record from the "category" table
	var product entity.Product
	result := db.First(&product)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	assert.Equal(t, false, product.IsActive)
}
