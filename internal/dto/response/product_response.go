package response

type ProductResponse struct {
	ID          uint64  `json:"id"`
	CategoryID  uint64  `json:"category_id"`
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}
