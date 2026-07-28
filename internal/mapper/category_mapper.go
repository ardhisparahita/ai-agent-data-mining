package mapper

import (
	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

func ToCategoryResponse(category *domain.Category) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func ToCategoryResponses(categories []domain.Category) []response.CategoryResponse {
	result := make([]response.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, *ToCategoryResponse(&category))
	}

	return result
}
