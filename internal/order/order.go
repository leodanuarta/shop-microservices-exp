// Package order is the Order Service: creates and lists orders.
// It validates items by calling product-service through ProductClient.
package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/example/shop-microservices/internal/httpx"
)

const StatusPending = "pending"

type Item struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Line struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}

// Order total is in whole rupiah (IDR).
type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Lines  []Line `json:"lines"`
	Total  int64  `json:"total"`
}

// Store is an in-memory order store, safe for concurrent use.
type Store struct {
	mu     sync.Mutex
	orders []Order
}

func NewStore() *Store { return &Store{} }

func (s *Store) Add(o Order) Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	o.ID = fmt.Sprintf("ord-%d", len(s.orders)+1)
	s.orders = append(s.orders, o)
	return o
}

// List returns all orders, or only those with the given status when non-empty.
func (s *Store) List(status string) []Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Order{}
	for _, o := range s.orders {
		if status == "" || o.Status == status {
			out = append(out, o)
		}
	}
	return out
}

// Get returns an order by ID.
func (s *Store) Get(id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, o := range s.orders {
		if o.ID == id {
			return o, true
		}
	}
	return Order{}, false
}

// Handler wires HTTP to the store and the product client.
type Handler struct {
	store    *Store
	products ProductClient
}

func NewHandler(s *Store, pc ProductClient) *Handler { return &Handler{store: s, products: pc} }

// Routes registers: GET /health, GET /orders, GET /orders/{id}, POST /orders.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpx.Health("order-service"))
	mux.HandleFunc("GET /orders", h.list)
	mux.HandleFunc("GET /orders/{id}", h.get)
	mux.HandleFunc("POST /orders", h.create)
	return mux
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, h.store.List(r.URL.Query().Get("status")))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	o, ok := h.store.Get(r.PathValue("id"))
	if !ok {
		httpx.Error(w, http.StatusNotFound, "order not found")
		return
	}
	httpx.JSON(w, http.StatusOK, o)
}

type createRequest struct {
	Items []Item `json:"items"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if len(req.Items) == 0 {
		httpx.Error(w, http.StatusBadRequest, "items must not be empty")
		return
	}

	var lines []Line
	var total int64
	for _, it := range req.Items {
		if it.ProductID == "" || it.Quantity <= 0 {
			httpx.Error(w, http.StatusBadRequest, "each item needs product_id and quantity > 0")
			return
		}
		p, err := h.products.GetProduct(r.Context(), it.ProductID)
		if errors.Is(err, ErrProductNotFound) {
			httpx.Error(w, http.StatusUnprocessableEntity, "unknown product: "+it.ProductID)
			return
		}
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, "product-service unavailable")
			return
		}
		if p.Stock < it.Quantity {
			httpx.Error(w, http.StatusConflict, "insufficient stock for "+p.ID)
			return
		}
		lines = append(lines, Line{ProductID: p.ID, Name: p.Name, Quantity: it.Quantity, UnitPrice: p.Price})
		total += p.Price * int64(it.Quantity)
	}

	created := h.store.Add(Order{Status: StatusPending, Lines: lines, Total: total})
	httpx.JSON(w, http.StatusCreated, created)
}
