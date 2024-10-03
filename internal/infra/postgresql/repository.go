package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/simpler-tha/internal/app"
)

type Repository struct {
	client *Client
}

func NewRepository(cl *Client) (Repository, error) {
	if cl == nil {
		return Repository{}, errors.New("client is nil")
	}

	return Repository{client: cl}, nil
}

func (r Repository) CreateProduct(ctx context.Context, p *app.Product) error {
	const sqlQuery = `
		INSERT INTO public.products (id, name, description, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.client.Conn.Exec(ctx, sqlQuery,
		p.ID, p.Name, p.Description, p.Price, p.CreatedAt.UTC(), p.UpdatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert product in the database: %w", err)
	}

	return nil
}

func (r Repository) UpdateProduct(ctx context.Context, p *app.Product) error {
	const sqlQuery = `
		UPDATE public.products
		SET name = $1, description = $2, price = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.client.Conn.Exec(ctx, sqlQuery,
		p.Name, p.Description, p.Price, p.UpdatedAt.UTC(), p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product with id %s in the database: %w", p.ID, err)
	}

	return nil
}

func (r Repository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	const sqlQuery = `
		DELETE FROM public.products
		WHERE id = $1
	`

	_, err := r.client.Conn.Exec(ctx, sqlQuery, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product with id %s from the database: %w", productID, err)
	}

	return nil
}

func (r Repository) GetProduct(ctx context.Context, productID uuid.UUID) (*app.Product, error) {
	const sqlQuery = `
		SELECT id, name, description, price, created_at, updated_at
		FROM public.products
		WHERE id = $1
	`

	var p app.Product

	err := r.client.Conn.QueryRow(ctx, sqlQuery, productID).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product with id %s from the database: %w", productID, err)
	}

	return &p, nil
}

func (r Repository) GetProducts(ctx context.Context, limit, offset int) ([]*app.Product, error) {
	const sqlQuery = `
		SELECT id, name, description, price, created_at, updated_at
		FROM public.products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.client.Conn.Query(ctx, sqlQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products from the database: %w", err)
	}
	defer rows.Close()

	var products []*app.Product

	for rows.Next() {
		var p app.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration over product rows: %w", err)
	}

	return products, nil
}
