package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	s, err := NewService(nil)
	assert.EqualError(t, err, "repository is nil")
	assert.Empty(t, s)
}

func TestService_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	dto := CreateProductDTO{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
	}

	tests := []struct {
		name                    string
		dto                     CreateProductDTO
		expRepoCreateProductErr error
		expErr                  error
	}{
		{
			name:                    "product was created successfully",
			dto:                     dto,
			expRepoCreateProductErr: nil,
			expErr:                  nil,
		},
		{
			name:                    "error creating product",
			dto:                     dto,
			expRepoCreateProductErr: errors.New("repo error"),
			expErr:                  fmt.Errorf("failed to create product: %w", errors.New("repo error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expectations
			mockRepository := NewMockrepository(ctrl)

			mockRepository.EXPECT().
				CreateProduct(gomock.Any(), gomock.Any()).
				Return(tt.expRepoCreateProductErr)

			// Exercise
			s, err := NewService(mockRepository)
			assert.NoError(t, err)

			p, err := s.CreateProduct(ctx, tt.dto)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p.ID)
				assert.Equal(t, tt.dto.Name, p.Name)
				assert.Equal(t, tt.dto.Description, p.Description)
				assert.Equal(t, tt.dto.Price, p.Price)
				assert.False(t, p.CreatedAt.IsZero())
				assert.False(t, p.UpdatedAt.IsZero())
			}
		})
	}
}

func TestService_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	productID := uuid.New()

	dto := UpdateProductDTO{
		ID:          productID,
		Name:        "New Product Name",
		Description: "New Product Description",
		Price:       200.0,
	}

	now := time.Now()

	expProduct := &Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name                    string
		dto                     UpdateProductDTO
		expRepoGetProductResult *Product
		expRepoGetProductErr    error
		expRepoUpdateProductErr error
		expErr                  error
	}{
		{
			name:                    "product was updated successfully",
			dto:                     dto,
			expRepoGetProductResult: expProduct,
			expRepoGetProductErr:    nil,
			expRepoUpdateProductErr: nil,
			expErr:                  nil,
		},
		{
			name:                    "error getting product",
			dto:                     dto,
			expRepoGetProductResult: nil,
			expRepoGetProductErr:    errors.New("repo error"),
			expRepoUpdateProductErr: nil,
			expErr:                  fmt.Errorf("failed to get product: %w", errors.New("repo error")),
		},
		{
			name:                    "error updating product",
			dto:                     dto,
			expRepoGetProductResult: expProduct,
			expRepoGetProductErr:    nil,
			expRepoUpdateProductErr: errors.New("repo error"),
			expErr:                  fmt.Errorf("failed to update product: %w", errors.New("repo error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expectations
			mockRepository := NewMockrepository(ctrl)

			mockRepository.EXPECT().
				GetProduct(gomock.Any(), productID).
				Return(tt.expRepoGetProductResult, tt.expRepoGetProductErr)

			if tt.expRepoGetProductErr == nil {
				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Return(tt.expRepoUpdateProductErr)
			}

			// Exercise
			s, err := NewService(mockRepository)
			assert.NoError(t, err)

			p, err := s.UpdateProduct(ctx, tt.dto)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, productID, p.ID)
				assert.Equal(t, tt.dto.Name, p.Name)
				assert.Equal(t, tt.dto.Description, p.Description)
				assert.Equal(t, tt.dto.Price, p.Price)
				assert.False(t, p.CreatedAt.IsZero())
				assert.False(t, p.UpdatedAt.IsZero())
			}
		})
	}
}

func TestService_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	productID := uuid.New()

	tests := []struct {
		name                    string
		expRepoDeleteProductErr error
		expErr                  error
	}{
		{
			name:                    "product was deleted successfully",
			expRepoDeleteProductErr: nil,
			expErr:                  nil,
		},
		{
			name:                    "error deleting product",
			expRepoDeleteProductErr: errors.New("repo error"),
			expErr:                  fmt.Errorf("failed to delete product: %w", errors.New("repo error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expectations
			mockRepository := NewMockrepository(ctrl)

			mockRepository.EXPECT().
				DeleteProduct(gomock.Any(), productID).
				Return(tt.expRepoDeleteProductErr)

			// Exercise
			s, err := NewService(mockRepository)
			assert.NoError(t, err)

			err = s.DeleteProduct(ctx, productID)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	productID := uuid.New()

	now := time.Now()

	expProduct := &Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name                    string
		expRepoGetProductResult *Product
		expRepoGetProductErr    error
		expErr                  error
	}{
		{
			name:                    "product was retrieved successfully",
			expRepoGetProductResult: expProduct,
			expRepoGetProductErr:    nil,
			expErr:                  nil,
		},
		{
			name:                    "error getting product",
			expRepoGetProductResult: nil,
			expRepoGetProductErr:    errors.New("repo error"),
			expErr:                  fmt.Errorf("failed to get product: %w", errors.New("repo error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expectations
			mockRepository := NewMockrepository(ctrl)

			mockRepository.EXPECT().
				GetProduct(gomock.Any(), productID).
				Return(tt.expRepoGetProductResult, tt.expRepoGetProductErr)

			// Exercise
			s, err := NewService(mockRepository)
			assert.NoError(t, err)

			p, err := s.GetProduct(ctx, productID)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expRepoGetProductResult.ID, p.ID)
				assert.Equal(t, tt.expRepoGetProductResult.Name, p.Name)
				assert.Equal(t, tt.expRepoGetProductResult.Description, p.Description)
				assert.Equal(t, tt.expRepoGetProductResult.Price, p.Price)
				assert.Equal(t, tt.expRepoGetProductResult.CreatedAt, p.CreatedAt)
				assert.Equal(t, tt.expRepoGetProductResult.UpdatedAt, p.UpdatedAt)
			}
		})
	}
}

func TestService_GetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	now := time.Now()

	expProductA := &Product{
		ID:          uuid.New(),
		Name:        "Test Product A",
		Description: "Test Description A",
		Price:       100.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	expProductB := &Product{
		ID:          uuid.New(),
		Name:        "Test Product B",
		Description: "Test Description B",
		Price:       200.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	limit := 20
	offset := 0

	tests := []struct {
		name                     string
		expRepoGetProductsResult []*Product
		expRepoGetProductsErr    error
		expErr                   error
	}{
		{
			name:                     "products were retrieved successfully",
			expRepoGetProductsResult: []*Product{expProductA, expProductB},
			expRepoGetProductsErr:    nil,
			expErr:                   nil,
		},
		{
			name:                     "error getting products",
			expRepoGetProductsResult: nil,
			expRepoGetProductsErr:    errors.New("repo error"),
			expErr:                   fmt.Errorf("failed to get products: %w", errors.New("repo error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expectations
			mockRepository := NewMockrepository(ctrl)

			mockRepository.EXPECT().
				GetProducts(gomock.Any(), limit, offset).
				Return(tt.expRepoGetProductsResult, tt.expRepoGetProductsErr)

			// Exercise
			s, err := NewService(mockRepository)
			assert.NoError(t, err)

			products, err := s.GetProducts(ctx, limit, offset)
			if tt.expErr != nil {
				assert.Equal(t, tt.expErr, err)
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Len(t, products, len(tt.expRepoGetProductsResult))

				productA := products[0]
				assert.Equal(t, tt.expRepoGetProductsResult[0].ID, productA.ID)
				assert.Equal(t, tt.expRepoGetProductsResult[0].Name, productA.Name)
				assert.Equal(t, tt.expRepoGetProductsResult[0].Description, productA.Description)
				assert.Equal(t, tt.expRepoGetProductsResult[0].Price, productA.Price)
				assert.Equal(t, tt.expRepoGetProductsResult[0].CreatedAt, productA.CreatedAt)
				assert.Equal(t, tt.expRepoGetProductsResult[0].UpdatedAt, productA.UpdatedAt)

				productB := products[1]
				assert.Equal(t, tt.expRepoGetProductsResult[1].ID, productB.ID)
				assert.Equal(t, tt.expRepoGetProductsResult[1].Name, productB.Name)
				assert.Equal(t, tt.expRepoGetProductsResult[1].Description, productB.Description)
				assert.Equal(t, tt.expRepoGetProductsResult[1].Price, productB.Price)
				assert.Equal(t, tt.expRepoGetProductsResult[1].CreatedAt, productB.CreatedAt)
				assert.Equal(t, tt.expRepoGetProductsResult[1].UpdatedAt, productB.UpdatedAt)
			}
		})
	}
}
