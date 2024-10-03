package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/simpler-tha/internal/app"
)

const (
	createProductEndpoint string = "POST /api/v1/products"
	updateProductEndpoint string = "PUT /api/v1/products/{product_id}"
	deleteProductEndpoint string = "DELETE /api/v1/products/{product_id}"
	getProductEndpoint    string = "GET /api/v1/products/{product_id}"
	getProductsEndpoint   string = "GET /api/v1/products"

	getProductsDefaultLimit  = 20
	getProductsDefaultOffset = 0
)

type Router struct {
	service service
}

type service interface {
	CreateProduct(ctx context.Context, dto app.CreateProductDTO) (*app.Product, error)
	UpdateProduct(ctx context.Context, dto app.UpdateProductDTO) (*app.Product, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	GetProduct(ctx context.Context, productID uuid.UUID) (*app.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*app.Product, error)
}

func NewRouter(s service) (Router, error) {
	if s == nil {
		return Router{}, errors.New("service cannot be nil")
	}
	return Router{service: s}, nil
}

func (r Router) RegisterRoutes() {
	http.HandleFunc(createProductEndpoint, r.createProductHandler)
	http.HandleFunc(updateProductEndpoint, r.updateProductHandler)
	http.HandleFunc(deleteProductEndpoint, r.deleteProductHandler)
	http.HandleFunc(getProductEndpoint, r.getProductHandler)
	http.HandleFunc(getProductsEndpoint, r.getProductsHandler)
}

type productRequestBody struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

func (r Router) createProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body productRequestBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dto := app.CreateProductDTO{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
	}

	p, err := r.service.CreateProduct(ctx, dto)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	productRes := newProductResponse(p)

	err = json.NewEncoder(w).Encode(productRes)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
	}
}

func (r Router) updateProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var body productRequestBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	productIDStr := req.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	dto := app.UpdateProductDTO{
		ID:          productID,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
	}

	p, err := r.service.UpdateProduct(ctx, dto)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	productRes := newProductResponse(p)

	err = json.NewEncoder(w).Encode(productRes)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
	}
}

func (r Router) deleteProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	productIDStr := req.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = r.service.DeleteProduct(ctx, productID)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (r Router) getProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	productIDStr := req.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := r.service.GetProduct(ctx, productID)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	productRes := newProductResponse(product)

	err = json.NewEncoder(w).Encode(productRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r Router) getProductsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	query := req.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = getProductsDefaultLimit
	}

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = getProductsDefaultOffset
	}

	products, err := r.service.GetProducts(ctx, limit, offset)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	var res productsResponse
	for _, p := range products {
		productRes := newProductResponse(p)
		res.Products = append(res.Products, productRes)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
	}
}

type productsResponse struct {
	Products []productResponse `json:"products"`
}

type productResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newProductResponse(p *app.Product) productResponse {
	return productResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
