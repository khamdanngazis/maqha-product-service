package service

import "fmt"

const (
	SuccessError   = 00
	SuccessMessage = "Success"

	//1 -99 general error
	GenaralSystemError        = 99
	GenaralSystemErrorMessage = "General System Error"

	//100 to 199: Authentication and authorization errors
	InvalidUsername        = 101
	InvalidUsernameMessage = "User not found"
	InvalidPassword        = 102
	InvalidPasswordMessage = "Invalid Password"

	//200 - 299 input validation error
	InvalidFormatError        = 201
	InvalidFormatErrorMessage = "Invalid Format Request"
	InvalidToken              = 202
	InvalidTokendMessage      = "Invalid Token"
	InvalidRequestError       = 203
	InvalidRequestMessage     = "Invalid Request %s"

	ProductNotFound            = 204
	ProductNotFoundMessage     = "Product Not Found"
	InvalidProductPrice        = 205
	InvalidProductPriceMessage = "Invalid Product Price"
	InvalidTotal               = 206
	InvalidTotalMessage        = "Invalid Total"

	//300 to 399: Database-related errors
	QueryError              = 301
	QueryErrorMessage       = "Error query database"
	UpdateQueryError        = 302
	UpdateQueryErrorMessage = "Error Update database"

	//600 to 699: Business-specific errors
	//product service error 600 -620
	DateCategoryNotFound        = 601
	DateCategoryNotFoundMessage = "Data Not Found"
)

// AppError represents an application-specific error.
type AppError struct {
	Code    int
	Message string
}

// NewAppError creates a new instance of AppError.
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewSuccessError() *AppError {
	return NewAppError(SuccessError, SuccessMessage)
}

func NewInvalidFormatError() *AppError {
	return NewAppError(InvalidFormatError, InvalidFormatErrorMessage)
}

func NewInvalidRequestError(s string) *AppError {
	return NewAppError(InvalidRequestError, fmt.Sprintf(InvalidRequestMessage, s))
}

func NewInvalidPasswordError() *AppError {
	return NewAppError(InvalidPassword, InvalidPasswordMessage)
}

func NewUserNotFoundError() *AppError {
	return NewAppError(InvalidUsername, InvalidUsernameMessage)
}

func NewInvalidTokenError() *AppError {
	return NewAppError(InvalidToken, InvalidTokendMessage)
}

func NewGeneralSystemError() *AppError {
	return NewAppError(GenaralSystemError, GenaralSystemErrorMessage)
}

func NewQueryDBError() *AppError {
	return NewAppError(QueryError, QueryErrorMessage)
}
func NewUpdateQueryDBError() *AppError {
	return NewAppError(UpdateQueryError, UpdateQueryErrorMessage)
}

func NewDateCategoryNotFoundError() *AppError {
	return NewAppError(DateCategoryNotFound, DateCategoryNotFoundMessage)
}

func NewProductNotFoundError() *AppError {
	return NewAppError(ProductNotFound, ProductNotFoundMessage)
}

func NewInvalidProductPriceError() *AppError {
	return NewAppError(InvalidProductPrice, InvalidProductPriceMessage)
}

func NewInvalidTotalError() *AppError {
	return NewAppError(InvalidTotal, InvalidTotalMessage)
}
