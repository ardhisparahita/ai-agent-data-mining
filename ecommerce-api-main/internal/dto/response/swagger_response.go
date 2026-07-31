package response

type ErrorSwaggerResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MessageSwaggerResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthSwaggerResponse struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    AuthResponse `json:"data"`
}

type RegisterSwaggerResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserSwaggerResponse struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
}

type CategorySwaggerResponse struct {
	Code    int              `json:"code"`
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}

type CategoryListSwaggerResponse struct {
	Code    int                `json:"code"`
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []CategoryResponse `json:"data"`
}

type ProductSwaggerResponse struct {
	Code    int             `json:"code"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    ProductResponse `json:"data"`
}

type ProductListSwaggerResponse struct {
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    ProductListResponse `json:"data"`
}

type AddressSwaggerResponse struct {
	Code    int             `json:"code"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    AddressResponse `json:"data"`
}

type AddressListSwaggerResponse struct {
	Code    int               `json:"code"`
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []AddressResponse `json:"data"`
}

type CartSwaggerResponse struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    CartResponse `json:"data"`
}

type CartListSwaggerResponse struct {
	Code    int            `json:"code"`
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []CartResponse `json:"data"`
}

type CheckoutSwaggerResponse struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    OrderResponse `json:"data"`
}

type OrderSwaggerResponse struct {
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    OrderDetailResponse `json:"data"`
}

type OrderListSwaggerResponse struct {
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []OrderListResponse `json:"data"`
}
