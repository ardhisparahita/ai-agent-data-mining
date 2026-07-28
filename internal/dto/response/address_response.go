package response

type AddressResponse struct {
	ID            uint64 `json:"id"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	PostalCode    string `json:"postal_code"`
}
