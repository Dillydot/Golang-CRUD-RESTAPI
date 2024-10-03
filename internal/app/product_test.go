package app

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	name := "Test Product"
	description := "Test Product Description"
	var price float32 = 100.0

	product := NewProduct(name, description, price)

	assert.NotNil(t, product.ID)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, description, product.Description)
	assert.Equal(t, price, product.Price)
	assert.False(t, product.CreatedAt.IsZero())
	assert.False(t, product.UpdatedAt.IsZero())
}

func TestProduct_Update(t *testing.T) {
	productID := uuid.New()
	now := time.Now()

	product := Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Product Description",
		Price:       200.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	newName := "Test Product 2"
	newDescription := "Test Product 2 Description"
	var newPrice float32 = 200.0

	product.Update(newName, newDescription, newPrice)

	assert.Equal(t, productID, product.ID)
	assert.Equal(t, newName, product.Name)
	assert.Equal(t, newDescription, product.Description)
	assert.Equal(t, newPrice, product.Price)
	assert.Equal(t, now, product.CreatedAt)
	assert.NotEqual(t, now, product.UpdatedAt)
}
