package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type repository interface {
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	GetProduct(ctx context.Context, productID uuid.UUID) (*Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*Product, error)
}

type Service struct {
	repository repository
}

func NewService(r repository) (Service, error) {
	if r == nil {
		return Service{}, errors.New("repository is nil")
	}

	return Service{
		repository: r,
	}, nil
}

func (s Service) CreateProduct(ctx context.Context, dto CreateProductDTO) (*Product, error) {
	p := NewProduct(dto.Name, dto.Description, dto.Price)

	err := s.repository.CreateProduct(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return p, nil
}

func (s Service) UpdateProduct(ctx context.Context, dto UpdateProductDTO) (*Product, error) {

	p, err := s.repository.GetProduct(ctx, dto.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	p.Update(dto.Name, dto.Description, dto.Price)

	err = s.repository.UpdateProduct(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return p, nil
}

func (s Service) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	err := s.repository.DeleteProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s Service) GetProduct(ctx context.Context, productID uuid.UUID) (*Product, error) {
	p, err := s.repository.GetProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return p, nil
}

func (s Service) GetProducts(ctx context.Context, limit, offset int) ([]*Product, error) {
	products, err := s.repository.GetProducts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

type CreateProductDTO struct {
	Name        string
	Description string
	Price       float32
}

type UpdateProductDTO struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float32
}
