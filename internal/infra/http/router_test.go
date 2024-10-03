package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/simpler-tha/internal/app"
)

func TestNewRouter(t *testing.T) {
	r, err := NewRouter(nil)
	assert.EqualError(t, err, "service cannot be nil")
	assert.Empty(t, r)
}

func TestRouter_createProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	reqBody := []byte(`{"name":"Test Product","description":"Test Description","price":100.0}`)

	createProductDTO := app.CreateProductDTO{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
	}

	productID, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd604e2c")

	now, _ := time.Parse(time.RFC3339, "2024-10-02T14:28:34Z")

	product := &app.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	responseBody := []byte("{\"id\":\"9f9f4340-6bf9-4948-808c-ebf2dd604e2c\",\"name\":\"Test Product\",\"description\":\"Test Description\",\"price\":100,\"created_at\":\"2024-10-02T14:28:34Z\",\"updated_at\":\"2024-10-02T14:28:34Z\"}\n")

	tests := []struct {
		name                          string
		reqBody                       []byte
		expServiceCreateProductResult *app.Product
		expServiceCreateProductError  error
		expStatus                     int
		expResponse                   []byte
	}{
		{
			name:                          "product created successfully",
			reqBody:                       reqBody,
			expServiceCreateProductResult: product,
			expServiceCreateProductError:  nil,
			expStatus:                     http.StatusCreated,
			expResponse:                   responseBody,
		},
		{
			name:                          "product could not be created",
			reqBody:                       reqBody,
			expServiceCreateProductResult: nil,
			expServiceCreateProductError:  errors.New("service error"),
			expStatus:                     http.StatusInternalServerError,
			expResponse:                   []byte("an error occurred\n"),
		},
		{
			name:                          "invalid request body",
			reqBody:                       []byte(`{`),
			expServiceCreateProductResult: nil,
			expServiceCreateProductError:  nil,
			expStatus:                     http.StatusBadRequest,
			expResponse:                   []byte("invalid request body\n"),
		},
	}

	for _, tt := range tests {
		mockService := NewMockservice(ctrl)

		if tt.expServiceCreateProductResult != nil || tt.expServiceCreateProductError != nil {
			mockService.
				EXPECT().
				CreateProduct(gomock.Any(), createProductDTO).
				Return(tt.expServiceCreateProductResult, tt.expServiceCreateProductError)
		}

		router, err := NewRouter(mockService)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(tt.reqBody))

		recorder := httptest.NewRecorder()

		router.createProductHandler(recorder, req)

		assert.Equal(t, tt.expStatus, recorder.Code)

		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)

		assert.Equal(t, tt.expResponse, b)
	}
}

func TestRouter_updateProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	reqBody := []byte(`{"name":"Test Product","description":"Test Description","price":100.0}`)

	productID, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd604e2c")

	updateProductDTO := app.UpdateProductDTO{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
	}

	now, _ := time.Parse(time.RFC3339, "2024-10-02T14:28:34Z")

	product := &app.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	responseBody := []byte("{\"id\":\"9f9f4340-6bf9-4948-808c-ebf2dd604e2c\",\"name\":\"Test Product\",\"description\":\"Test Description\",\"price\":100,\"created_at\":\"2024-10-02T14:28:34Z\",\"updated_at\":\"2024-10-02T14:28:34Z\"}\n")

	tests := []struct {
		name                          string
		productID                     string
		reqBody                       []byte
		expServiceUpdateProductResult *app.Product
		expServiceUpdateProductError  error
		expStatus                     int
		expResponse                   []byte
	}{
		{
			name:                          "product updated successfully",
			productID:                     productID.String(),
			reqBody:                       reqBody,
			expServiceUpdateProductResult: product,
			expServiceUpdateProductError:  nil,
			expStatus:                     http.StatusOK,
			expResponse:                   responseBody,
		},
		{
			name:                          "product could not be updated",
			productID:                     productID.String(),
			reqBody:                       reqBody,
			expServiceUpdateProductResult: nil,
			expServiceUpdateProductError:  errors.New("service error"),
			expStatus:                     http.StatusInternalServerError,
			expResponse:                   []byte("an error occurred\n"),
		},
		{
			name:                          "invalid request body",
			productID:                     productID.String(),
			reqBody:                       []byte(`{`),
			expServiceUpdateProductResult: nil,
			expServiceUpdateProductError:  nil,
			expStatus:                     http.StatusBadRequest,
			expResponse:                   []byte("invalid request body\n"),
		},
		{
			name:                          "invalid product ID",
			productID:                     "",
			reqBody:                       reqBody,
			expServiceUpdateProductResult: nil,
			expServiceUpdateProductError:  nil,
			expStatus:                     http.StatusBadRequest,
			expResponse:                   []byte("invalid product id\n"),
		},
	}

	for _, tt := range tests {
		mockService := NewMockservice(ctrl)

		if tt.expServiceUpdateProductResult != nil || tt.expServiceUpdateProductError != nil {
			mockService.
				EXPECT().
				UpdateProduct(gomock.Any(), updateProductDTO).
				Return(tt.expServiceUpdateProductResult, tt.expServiceUpdateProductError)
		}

		router, err := NewRouter(mockService)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/products", bytes.NewBuffer(tt.reqBody))
		if tt.productID != "" {
			req.SetPathValue("product_id", tt.productID)
		}

		recorder := httptest.NewRecorder()

		router.updateProductHandler(recorder, req)

		assert.Equal(t, tt.expStatus, recorder.Code)

		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)

		assert.Equal(t, tt.expResponse, b)
	}
}

func TestRouter_deleteProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	productID, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd604e2c")

	tests := []struct {
		name                         string
		productID                    string
		expServiceDeleteProductError error
		expStatus                    int
		expResponse                  []byte
	}{
		{
			name:                         "product deleted successfully",
			productID:                    productID.String(),
			expServiceDeleteProductError: nil,
			expStatus:                    http.StatusNoContent,
			expResponse:                  []byte{},
		},
		{
			name:                         "product could not be deleted",
			productID:                    productID.String(),
			expServiceDeleteProductError: errors.New("service error"),
			expStatus:                    http.StatusInternalServerError,
			expResponse:                  []byte("an error occurred\n"),
		},
		{
			name:                         "invalid product ID",
			productID:                    "",
			expServiceDeleteProductError: nil,
			expStatus:                    http.StatusBadRequest,
			expResponse:                  []byte("invalid product id\n"),
		},
	}

	for _, tt := range tests {
		mockService := NewMockservice(ctrl)

		if tt.productID != "" {
			mockService.
				EXPECT().
				DeleteProduct(gomock.Any(), productID).
				Return(tt.expServiceDeleteProductError)
		}

		router, err := NewRouter(mockService)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/{product_id}", nil)
		if tt.productID != "" {
			req.SetPathValue("product_id", tt.productID)
		}

		recorder := httptest.NewRecorder()

		router.deleteProductHandler(recorder, req)

		assert.Equal(t, tt.expStatus, recorder.Code)

		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)

		assert.Equal(t, tt.expResponse, b)
	}
}

func TestRouter_getProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	productID, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd604e2c")

	now, _ := time.Parse(time.RFC3339, "2024-10-02T14:28:34Z")

	product := &app.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	responseBody := []byte("{\"id\":\"9f9f4340-6bf9-4948-808c-ebf2dd604e2c\",\"name\":\"Test Product\",\"description\":\"Test Description\",\"price\":100,\"created_at\":\"2024-10-02T14:28:34Z\",\"updated_at\":\"2024-10-02T14:28:34Z\"}\n")

	tests := []struct {
		name                       string
		productID                  string
		expServiceGetProductResult *app.Product
		expServiceGetProductError  error
		expStatus                  int
		expResponse                []byte
	}{
		{
			name:                       "product retrieved successfully",
			productID:                  productID.String(),
			expServiceGetProductResult: product,
			expServiceGetProductError:  nil,
			expStatus:                  http.StatusOK,
			expResponse:                responseBody,
		},
		{
			name:                       "product could not be retrieved",
			productID:                  productID.String(),
			expServiceGetProductResult: nil,
			expServiceGetProductError:  errors.New("service error"),
			expStatus:                  http.StatusInternalServerError,
			expResponse:                []byte("an error occurred\n"),
		},
		{
			name:                      "invalid product ID",
			productID:                 "",
			expServiceGetProductError: nil,
			expStatus:                 http.StatusBadRequest,
			expResponse:               []byte("invalid product id\n"),
		},
	}

	for _, tt := range tests {
		mockService := NewMockservice(ctrl)

		if tt.productID != "" {
			mockService.
				EXPECT().
				GetProduct(gomock.Any(), productID).
				Return(tt.expServiceGetProductResult, tt.expServiceGetProductError)
		}

		router, err := NewRouter(mockService)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/{product_id}", nil)
		if tt.productID != "" {
			req.SetPathValue("product_id", tt.productID)
		}

		recorder := httptest.NewRecorder()

		router.getProductHandler(recorder, req)

		assert.Equal(t, tt.expStatus, recorder.Code)

		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)

		assert.Equal(t, tt.expResponse, b)
	}
}

func TestRouter_getProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	productIDA, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd604e2c")
	productIDB, _ := uuid.Parse("9f9f4340-6bf9-4948-808c-ebf2dd60303d")

	now, _ := time.Parse(time.RFC3339, "2024-10-02T14:28:34Z")

	productA := &app.Product{
		ID:          productIDA,
		Name:        "Test Product A",
		Description: "Test Description A",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	productB := &app.Product{
		ID:          productIDB,
		Name:        "Test Product B",
		Description: "Test Description B",
		Price:       200.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	responseBody := []byte("{\"products\":[" +
		"{\"id\":\"9f9f4340-6bf9-4948-808c-ebf2dd604e2c\",\"name\":\"Test Product A\",\"description\":\"Test Description A\",\"price\":100,\"created_at\":\"2024-10-02T14:28:34Z\",\"updated_at\":\"2024-10-02T14:28:34Z\"}," +
		"{\"id\":\"9f9f4340-6bf9-4948-808c-ebf2dd60303d\",\"name\":\"Test Product B\",\"description\":\"Test Description B\",\"price\":200,\"created_at\":\"2024-10-02T14:28:34Z\",\"updated_at\":\"2024-10-02T14:28:34Z\"}" +
		"]}\n")

	tests := []struct {
		name                        string
		limit                       string
		offset                      string
		expServiceGetProductsResult []*app.Product
		expServiceGetProductsError  error
		expStatus                   int
		expResponse                 []byte
	}{
		{
			name:                        "product retrieved successfully",
			limit:                       "10",
			offset:                      "0",
			expServiceGetProductsResult: []*app.Product{productA, productB},
			expServiceGetProductsError:  nil,
			expStatus:                   http.StatusOK,
			expResponse:                 responseBody,
		},
		{
			name:                        "product retrieved successfully with default limit and offset",
			limit:                       "",
			offset:                      "",
			expServiceGetProductsResult: []*app.Product{productA, productB},
			expServiceGetProductsError:  nil,
			expStatus:                   http.StatusOK,
			expResponse:                 responseBody,
		},
		{
			name:                        "product could not be retrieved",
			limit:                       "10",
			offset:                      "0",
			expServiceGetProductsResult: nil,
			expServiceGetProductsError:  errors.New("service error"),
			expStatus:                   http.StatusInternalServerError,
			expResponse:                 []byte("an error occurred\n"),
		},
	}

	for _, tt := range tests {
		mockService := NewMockservice(ctrl)

		limit := tt.limit
		if limit == "" {
			limit = "5"
		}
		limitNumber, _ := strconv.Atoi(limit)

		offset := tt.offset
		if offset == "" {
			offset = "0"
		}
		offsetNumber, _ := strconv.Atoi(offset)

		mockService.
			EXPECT().
			GetProducts(gomock.Any(), limitNumber, offsetNumber).
			Return(tt.expServiceGetProductsResult, tt.expServiceGetProductsError)

		router, err := NewRouter(mockService)
		assert.NoError(t, err)
		assert.NotNil(t, router)

		url := fmt.Sprintf("/api/v1/products?limit=%s&offset=%s", tt.limit, tt.offset)
		req := httptest.NewRequest(http.MethodGet, url, nil)

		recorder := httptest.NewRecorder()

		router.getProductsHandler(recorder, req)

		assert.Equal(t, tt.expStatus, recorder.Code)

		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)

		assert.Equal(t, tt.expResponse, b)
	}
}
