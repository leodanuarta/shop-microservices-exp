package order

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeProducts map[string]Product

func (f fakeProducts) GetProduct(_ context.Context, id string) (Product, error) {
	p, ok := f[id]
	if !ok {
		return Product{}, ErrProductNotFound
	}
	return p, nil
}

func newHandler() *Handler {
	return NewHandler(NewStore(), fakeProducts{
		"p-1": {ID: "p-1", Name: "Kaos", Price: 79000, Stock: 5},
		"p-2": {ID: "p-2", Name: "Sepatu", Price: 450000, Stock: 1},
	})
}

func post(h *Handler, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body)))
	return rec
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler().Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name, body string
		wantStatus int
	}{
		{"ok", `{"items":[{"product_id":"p-1","quantity":2}]}`, http.StatusCreated},
		{"bad json", `{`, http.StatusBadRequest},
		{"empty items", `{"items":[]}`, http.StatusBadRequest},
		{"zero quantity", `{"items":[{"product_id":"p-1","quantity":0}]}`, http.StatusBadRequest},
		{"unknown product", `{"items":[{"product_id":"zzz","quantity":1}]}`, http.StatusUnprocessableEntity},
		{"insufficient stock", `{"items":[{"product_id":"p-2","quantity":2}]}`, http.StatusConflict},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if rec := post(newHandler(), tc.body); rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tc.wantStatus, rec.Body)
			}
		})
	}
}

func TestCreateOrderComputesTotal(t *testing.T) {
	rec := post(newHandler(), `{"items":[{"product_id":"p-1","quantity":2},{"product_id":"p-2","quantity":1}]}`)
	var o Order
	if err := json.Unmarshal(rec.Body.Bytes(), &o); err != nil {
		t.Fatal(err)
	}
	if want := int64(2*79000 + 450000); o.Total != want {
		t.Fatalf("total = %d, want %d", o.Total, want)
	}
	if o.Status != StatusPending || o.ID != "ord-1" {
		t.Fatalf("unexpected order: %+v", o)
	}
}

func TestListOrdersFilter(t *testing.T) {
	h := newHandler()
	post(h, `{"items":[{"product_id":"p-1","quantity":1}]}`)

	tests := []struct {
		path      string
		wantCount int
	}{
		{"/orders", 1},
		{"/orders?status=pending", 1},
		{"/orders?status=paid", 0},
	}
	for _, tc := range tests {
		rec := httptest.NewRecorder()
		h.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		var got []Order
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got == nil || len(got) != tc.wantCount {
			t.Fatalf("%s: got %v, want len %d", tc.path, got, tc.wantCount)
		}
	}
}
