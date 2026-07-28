package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
)

func setupProductHandler() (*fiber.App, *MockProductService) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Post("/products", handler.Create)

	return app, service
}

func TestProductCreateSuccess(t *testing.T) {

	app, service := setupProductHandler()

	reqBody := request.CreateProductRequest{
		Name:        "Mechanical Keyboard",
		Description: "Gaming Keyboard",
		Price:       750000,
		Stock:       10,
		CategoryID:  1,
	}

	service.On(
		"Create",
		mock.Anything,
		reqBody,
	).Return(
		&response.ProductResponse{
			ID:          1,
			Name:        "Mechanical Keyboard",
			Description: "Gaming Keyboard",
			Price:       750000,
			Stock:       10,
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
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

func TestProductCreateValidationError(t *testing.T) {

	app, _ := setupProductHandler()

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
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

func TestProductCreateServiceError(t *testing.T) {

	app, service := setupProductHandler()

	reqBody := request.CreateProductRequest{
		Name:        "Mechanical Keyboard",
		Description: "Gaming Keyboard",
		Price:       750000,
		Stock:       10,
		CategoryID:  1,
	}

	service.On(
		"Create",
		mock.Anything,
		reqBody,
	).Return(
		nil,
		utils.Conflict("product already exists"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
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
		fiber.StatusConflict,
		resp.StatusCode,
	)

	service.AssertExpectations(t)
}

func TestProductFindAllSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Get("/products", handler.FindAll)

	reqQuery := request.ProductQueryRequest{
		Page:       1,
		Limit:      10,
		Search:     "",
		CategoryID: 0,
	}

	service.On(
		"FindAll",
		mock.Anything,
		reqQuery,
	).Return(
		&response.ProductListResponse{
			Items: []response.ProductResponse{
				{
					ID:          1,
					Name:        "Mechanical Keyboard",
					Description: "Gaming Keyboard",
					Price:       750000,
					Stock:       10,
				},
			},
			Page:       1,
			Limit:      10,
			TotalRows:  1,
			TotalPages: 1,
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products?page=1&limit=10",
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

func TestProductFindAllServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Get("/products", handler.FindAll)

	reqQuery := request.ProductQueryRequest{
		Page:       1,
		Limit:      10,
		Search:     "",
		CategoryID: 0,
	}

	service.On(
		"FindAll",
		mock.Anything,
		reqQuery,
	).Return(
		nil,
		utils.BadRequest("failed to get products"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products?page=1&limit=10",
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

func TestProductFindByIDSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Get("/products/:id", handler.FindByID)

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		&response.ProductResponse{
			ID:          1,
			Name:        "Mechanical Keyboard",
			Description: "Gaming Keyboard",
			Price:       750000,
			Stock:       10,
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/1",
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

func TestProductFindByIDInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Get("/products/:id", handler.FindByID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/abc",
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

func TestProductFindByIDServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Get("/products/:id", handler.FindByID)

	service.On(
		"FindByID",
		mock.Anything,
		uint64(1),
	).Return(
		nil,
		utils.NotFound("product not found"),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/1",
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

func TestProductUpdateSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Put("/products/:id", handler.Update)

	reqBody := request.UpdateProductRequest{
		Name:        "Mechanical Keyboard RGB",
		Description: "Gaming Keyboard RGB",
		Price:       850000,
		Stock:       20,
		CategoryID:  1,
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		&response.ProductResponse{
			ID:          1,
			Name:        "Mechanical Keyboard RGB",
			Description: "Gaming Keyboard RGB",
			Price:       850000,
			Stock:       20,
		},
		nil,
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
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

func TestProductUpdateInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Put("/products/:id", handler.Update)

	reqBody := request.UpdateProductRequest{
		Name:        "Mechanical Keyboard RGB",
		Description: "Gaming Keyboard RGB",
		Price:       850000,
		Stock:       20,
		CategoryID:  1,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/abc",
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

func TestProductUpdateValidationError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Put("/products/:id", handler.Update)

	reqBody := map[string]any{
		"name": "",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
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

func TestProductUpdateServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Put("/products/:id", handler.Update)

	reqBody := request.UpdateProductRequest{
		Name:        "Mechanical Keyboard RGB",
		Description: "Gaming Keyboard RGB",
		Price:       850000,
		Stock:       20,
		CategoryID:  1,
	}

	service.On(
		"Update",
		mock.Anything,
		uint64(1),
		reqBody,
	).Return(
		nil,
		utils.NotFound("product not found"),
	)

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
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

func TestProductDeleteSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Delete("/products/:id", handler.Delete)

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/1",
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

func TestProductDeleteInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Delete("/products/:id", handler.Delete)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/abc",
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

func TestProductDeleteServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Delete("/products/:id", handler.Delete)

	service.On(
		"Delete",
		mock.Anything,
		uint64(1),
	).Return(
		utils.NotFound("product not found"),
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/1",
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

func TestProductUploadImageSuccess(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Post("/products/:id/image", handler.UploadImage)

	// buat folder upload agar SaveFile tidak gagal
	err := os.MkdirAll("./uploads/products", os.ModePerm)
	assert.NoError(t, err)

	defer os.RemoveAll("./uploads")

	// buat temporary image
	tmpFile, err := os.CreateTemp("", "*.jpg")
	assert.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("dummy image"))
	assert.NoError(t, err)

	tmpFile.Close()

	// multipart body
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "image.jpg")
	assert.NoError(t, err)

	file, err := os.Open(tmpFile.Name())
	assert.NoError(t, err)

	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	writer.Close()

	service.On(
		"UploadImage",
		mock.Anything,
		uint64(1),
		mock.Anything,
	).Return(
		&response.ProductResponse{
			ID:          1,
			Name:        "Keyboard",
			Description: "Gaming Keyboard",
			ImageURL:    "/uploads/products/image.jpg",
		},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products/1/image",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
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

func TestProductUploadImageInvalidID(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Post("/products/:id/image", handler.UploadImage)

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/products/abc/image",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestProductUploadImageNoFile(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Post("/products/:id/image", handler.UploadImage)

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/products/1/image",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(
		t,
		fiber.StatusBadRequest,
		resp.StatusCode,
	)
}

func TestProductUploadImageServiceError(t *testing.T) {

	app := fiber.New()

	service := new(MockProductService)

	handler := ProductHandler{
		Service: service,
	}

	app.Post("/products/:id/image", handler.UploadImage)

	err := os.MkdirAll("./uploads/products", os.ModePerm)
	assert.NoError(t, err)

	defer os.RemoveAll("./uploads")

	tmpFile, err := os.CreateTemp("", "*.jpg")
	assert.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("dummy image"))
	assert.NoError(t, err)

	tmpFile.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "image.jpg")
	assert.NoError(t, err)

	file, err := os.Open(tmpFile.Name())
	assert.NoError(t, err)

	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	writer.Close()

	service.On(
		"UploadImage",
		mock.Anything,
		uint64(1),
		mock.Anything,
	).Return(
		nil,
		utils.NotFound("product not found"),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products/1/image",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
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