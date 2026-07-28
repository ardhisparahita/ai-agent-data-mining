package mapper

import (
	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

func ToCartResponse(cart *domain.Cart) *response.CartResponse {
	subtotal := cart.Product.Price * float64(cart.Quantity)

	return &response.CartResponse{
		ID:       cart.Product.ID,
		Quantity: cart.Quantity,
		Subtotal: subtotal,
		Product: response.CartProductResponse{
			ID:       cart.Product.ID,
			Name:     cart.Product.Name,
			Price:    cart.Product.Price,
			ImageURL: cart.Product.ImageURL,
		},
	}
}

func ToCartResponses(carts []domain.Cart) ([]response.CartResponse, float64) {

	result := make([]response.CartResponse, 0, len(carts))
	var grandTotal float64

	for _, cart := range carts {
		subtotal := cart.Product.Price * float64(cart.Quantity)

		grandTotal += subtotal

		result = append(result, *ToCartResponse(&cart))
	}

	return result, grandTotal
}
