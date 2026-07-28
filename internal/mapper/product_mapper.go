package mapper

import (
	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

func ToProductResponse(product *domain.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Category:    product.Category.Name,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
	}
}

func ToProductResponses(products []domain.Product) []response.ProductResponse {
	result := make([]response.ProductResponse, 0, len(products))

	for _, product := range products {
		result = append(result, *ToProductResponse(&product))
	}

	return result
}
