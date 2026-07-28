package request

type CheckoutRequest struct {
	AddressID     uint64 `json:"address_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}
