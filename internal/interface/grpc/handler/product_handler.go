package handler

import (
	"context"
	"maqhaa/product_service/internal/app/service"
	pb "maqhaa/product_service/internal/interface/grpc/model" // Update with your actual package name
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductGRPCHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}
func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, appError := h.productService.GetProductByID(ctx, uint(req.ProductId), req.Token)
	var response *pb.GetProductResponse

	if appError.Code != service.SuccessError {
		response = &pb.GetProductResponse{
			Code:    int32(appError.Code),
			Message: appError.Message,
			Data:    nil,
		}
		return response, nil
	}

	response = &pb.GetProductResponse{
		Code:    int32(appError.Code),
		Message: appError.Message,
		Data: &pb.ProductData{
			Id:          uint32(product.ID),
			CategoryId:  uint32(product.CategoryID),
			Name:        product.Name,
			Price:       float32(product.Price),
			Description: product.Description,
			Image:       product.Image,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	return response, nil
}
