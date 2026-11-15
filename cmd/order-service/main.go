package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/shop-microservices/internal/order"
)

func main() {
	addr := ":" + getenv("PORT", "8082")
	productURL := getenv("PRODUCT_SERVICE_URL", "http://localhost:8081")

	h := order.NewHandler(order.NewStore(), order.NewHTTPProductClient(productURL))
	log.Printf("order-service listening on %s (product-service: %s)", addr, productURL)
	log.Fatal(http.ListenAndServe(addr, h.Routes()))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
