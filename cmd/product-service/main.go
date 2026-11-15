package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/shop-microservices/internal/product"
)

func main() {
	addr := ":" + getenv("PORT", "8081")
	h := product.NewHandler(product.NewStore())
	log.Printf("product-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, h.Routes()))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
