package product

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func do(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	NewHandler(NewStore()).Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealth(t *testing.T) {
	rec := do(t, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestListProducts(t *testing.T) {
	tests := []struct {
		name, path string
		wantCount  int
	}{
		{"all", "/products", 5},
		{"filter category", "/products?category=electronics", 2},
		{"case insensitive", "/products?category=FASHION", 2},
		{"no match -> empty", "/products?category=toys", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, tc.path)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("content-type = %q", ct)
			}
			var got []Product
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got == nil || len(got) != tc.wantCount {
				t.Fatalf("got %v (len %d), want len %d and non-nil", got, len(got), tc.wantCount)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	if rec := do(t, "/products/p-3"); rec.Code != http.StatusOK {
		t.Fatalf("existing: status = %d, want 200", rec.Code)
	}
	if rec := do(t, "/products/nope"); rec.Code != http.StatusNotFound {
		t.Fatalf("missing: status = %d, want 404", rec.Code)
	}
}
