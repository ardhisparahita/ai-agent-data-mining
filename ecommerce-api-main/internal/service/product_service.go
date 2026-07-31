package service

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

type ProductService interface {
	Create(ctx context.Context, req request.CreateProductRequest) (*response.ProductResponse, error)
	FindAll(ctx context.Context, req request.ProductQueryRequest) (*response.ProductListResponse, error)
	FindByID(ctx context.Context, id uint64) (*response.ProductResponse, error)
	Update(ctx context.Context, id uint64, req request.UpdateProductRequest) (*response.ProductResponse, error)
	Delete(ctx context.Context, id uint64) error
	UploadImage(ctx context.Context, id uint64, imageURL string) (*response.ProductResponse, error)
}
