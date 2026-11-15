.PHONY: test run-product run-order
test:
	go test ./...
run-product:
	PORT=8081 go run ./cmd/product-service
run-order:
	PORT=8082 PRODUCT_SERVICE_URL=http://localhost:8081 go run ./cmd/order-service
