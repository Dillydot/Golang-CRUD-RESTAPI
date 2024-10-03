package app

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(name, description string, price float32) *Product {
	now := time.Now().UTC()

	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) Update(name, description string, price float32) {
	p.Name = name
	p.Description = description
	p.Price = price
	p.UpdatedAt = time.Now().UTC()
}
