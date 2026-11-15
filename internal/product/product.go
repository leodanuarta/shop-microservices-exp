// Package product is the Product Service: catalog data + HTTP handlers.
package product

import (
	"net/http"
	"strings"

	"github.com/example/shop-microservices/internal/httpx"
)

// Product is a catalog item. Price is in whole rupiah (IDR) to avoid float money.
type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int64  `json:"price"`
	Stock    int    `json:"stock"`
}

// Store is an in-memory catalog (no database needed for the workshop).
type Store struct {
	items []Product
}

// NewStore returns a store seeded with dummy data.
func NewStore() *Store {
	return &Store{items: []Product{
		{ID: "p-1", Name: "Kaos Polos Hitam", Category: "fashion", Price: 79000, Stock: 50},
		{ID: "p-2", Name: "Sepatu Lari", Category: "fashion", Price: 450000, Stock: 10},
		{ID: "p-3", Name: "Earphone Bluetooth", Category: "electronics", Price: 199000, Stock: 25},
		{ID: "p-4", Name: "Power Bank 10000mAh", Category: "electronics", Price: 159000, Stock: 0},
		{ID: "p-5", Name: "Kopi Arabika 250g", Category: "food", Price: 65000, Stock: 100},
	}}
}

// List returns all products, or only those in category when it is non-empty.
// Category matching is case-insensitive. No match returns an empty (non-nil) slice.
func (s *Store) List(category string) []Product {
	out := []Product{}
	for _, p := range s.items {
		if category == "" || strings.EqualFold(p.Category, category) {
			out = append(out, p)
		}
	}
	return out
}

// Get returns a product by ID.
func (s *Store) Get(id string) (Product, bool) {
	for _, p := range s.items {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

// Handler exposes the Store over HTTP.
type Handler struct{ store *Store }

func NewHandler(s *Store) *Handler { return &Handler{store: s} }

// Routes registers: GET /health, GET /products, GET /products/{id}.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpx.Health("product-service"))
	mux.HandleFunc("GET /products", h.list)
	mux.HandleFunc("GET /products/{id}", h.get)
	return mux
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	httpx.JSON(w, http.StatusOK, h.store.List(category))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	p, ok := h.store.Get(r.PathValue("id"))
	if !ok {
		httpx.Error(w, http.StatusNotFound, "product not found")
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}
