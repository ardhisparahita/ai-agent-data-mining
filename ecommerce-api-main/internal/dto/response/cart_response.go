package response

type CartResponse struct {
	ID       uint64              `json:"id"`
	Quantity int                 `json:"quantity"`
	Subtotal float64             `json:"subtotal"`
	Product  CartProductResponse `json:"product"`
}

type CartProductResponse struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
}

type CartListResponse struct {
	Items      []CartResponse `json:"items"`
	TotalItems int            `json:"total_item"`
	GrandTotal float64        `json:"grand_total"`
}
