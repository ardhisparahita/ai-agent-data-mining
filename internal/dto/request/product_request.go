package request

type CreateProductRequest struct {
	CategoryID  uint64  `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	CategoryID  uint64  `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}
