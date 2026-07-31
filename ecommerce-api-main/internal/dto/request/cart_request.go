package request

type AddToCartRequest struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}
type UpdateCartRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}
