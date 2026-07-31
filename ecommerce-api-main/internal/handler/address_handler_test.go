package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
)

func TestAddressCreateSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Post("/addresses", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Create(c)
	})

	reqBody := request.CreateAddressRequest{
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	service.On(
		"Create",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.AddressResponse{
			ID:            1,
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/addresses",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusCreated,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressCreateValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Post("/addresses", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Create(c)
	})

	reqBody := map[string]any{
		"recipient_name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/addresses",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestAddressCreateServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Post("/addresses", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Create(c)
	})

	reqBody := request.CreateAddressRequest{
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	service.On(
		"Create",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.BadRequest("failed create address"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/addresses",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressFindAllSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Get("/addresses", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		[]response.AddressResponse{
			{
				ID:            1,
				RecipientName: "John Doe",
				Phone:         "08123456789",
				Address:       "Jl. Malioboro",
				City:          "Yogyakarta",
				PostalCode:    "55213",
			},
			{
				ID:            2,
				RecipientName: "Jane Doe",
				Phone:         "08129876543",
				Address:       "Jl. Kaliurang",
				City:          "Yogyakarta",
				PostalCode:    "55281",
			},
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/addresses",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressFindAllServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Get("/addresses", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindAll(c)
	})

	service.On(
		"FindAllByUserID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.BadRequest("failed get addresses"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/addresses",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressFindByIDSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Get("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindByID(c)
	})

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		&response.AddressResponse{
			ID:            1,
			RecipientName: "John Doe",
			Phone:         "08123456789",
			Address:       "Jl. Malioboro",
			City:          "Yogyakarta",
			PostalCode:    "55213",
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/addresses/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressFindByIDInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Get("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindByID(c)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/addresses/abc",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestAddressFindByIDServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Get("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.FindByID(c)
	})

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		nil,
		utils.NotFound("address not found"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/addresses/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusNotFound,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressUpdateSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Put("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateAddressRequest{
		RecipientName: "John Doe Updated",
		Phone:         "08123456780",
		Address:       "Jl. Gejayan",
		City:          "Yogyakarta",
		PostalCode:    "55281",
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		uint64(1),
		reqBody,
	).Return(
		&response.AddressResponse{
			ID:            1,
			RecipientName: "John Doe Updated",
			Phone:         "08123456780",
			Address:       "Jl. Gejayan",
			City:          "Yogyakarta",
			PostalCode:    "55281",
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/addresses/1",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressUpdateInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Put("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateAddressRequest{
		RecipientName: "John Doe",
		Phone:         "08123456789",
		Address:       "Jl. Malioboro",
		City:          "Yogyakarta",
		PostalCode:    "55213",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/addresses/abc",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestAddressUpdateValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Put("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := map[string]any{
		"recipient_name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/addresses/1",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestAddressUpdateServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Put("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Update(c)
	})

	reqBody := request.UpdateAddressRequest{
		RecipientName: "John Doe Updated",
		Phone:         "08123456780",
		Address:       "Jl. Gejayan",
		City:          "Yogyakarta",
		PostalCode:    "55281",
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.NotFound("address not found"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/addresses/1",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusNotFound,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressDeleteSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Delete("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/addresses/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusOK,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestAddressDeleteInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Delete("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/addresses/abc",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestAddressDeleteServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockAddressService)

	handler := AddressHandler{
		Service: service,
	}

	app.Delete("/addresses/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uint64(1))
		return handler.Delete(c)
	})

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
		uint64(1),
	).Return(
		utils.NotFound("address not found"),
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/addresses/1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusNotFound,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}
