package mapper

import (
	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
)

func ToAddressResponse(address *domain.Address) *response.AddressResponse {
	return &response.AddressResponse{
		ID:            address.ID,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Address:       address.Address,
		City:          address.City,
		PostalCode:    address.PostalCode,
	}
}

func ToAddressResponses(addresses []domain.Address) []response.AddressResponse {
	result := make([]response.AddressResponse, 0, len(addresses))

	for _, address := range addresses {
		result = append(result, *ToAddressResponse(&address))
	}

	return result
}
