package request

type ProductQueryRequest struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	CategoryID uint64 `query:"category_id"`

	MinPrice float64 `query:"min_price"`
	MaxPrice float64 `query:"max_price"`

	SortBy string `query:"sort_by"`
	Order  string `query:"order"`
}
