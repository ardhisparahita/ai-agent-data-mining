package request

type CreateAddressRequest struct {
	RecipientName string `json:"recipient_name" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	PostalCode    string `json:"postal_code" validate:"required"`
}

type UpdateAddressRequest struct {
	RecipientName string `json:"recipient_name" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	PostalCode    string `json:"postal_code" validate:"required"`
}
