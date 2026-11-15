# shop-microservices

Dummy study case for the workshop **"From Natural Language to Executable Code"**.
Two tiny Go microservices for an e-commerce domain. Standard library only (Go 1.22+), in-memory data.

```
client --> order-service :8082 --HTTP--> product-service :8081
```

| Service | Endpoints |
|---|---|
| product-service | `GET /health`, `GET /products?category=`, `GET /products/{id}` |
| order-service | `GET /health`, `GET /orders?status=`, `POST /orders` |

## Run
```bash
make run-product   # terminal 1
make run-order     # terminal 2
make test          # go test ./...
```

## Try it
```bash
curl localhost:8081/products?category=electronics
curl -X POST localhost:8082/orders -d '{"items":[{"product_id":"p-1","quantity":2}]}'
curl "localhost:8082/orders?status=pending"
```
Error cases: `p-4` has stock 0 (409), unknown id (422), empty items (400).

## Layout
```
cmd/<service>/main.go     entry point (PORT, PRODUCT_SERVICE_URL env)
internal/product          catalog store + handlers + tests
internal/order            order store + handlers + ProductClient + tests
internal/httpx            JSON/error/health helpers
docs/SPEC.md              acceptance criteria (spec-first)
docs/PROMPTS.md           prompt template for the agent
docs/CHALLENGES.md        ownership challenges A/B/C
.agents/skills/go-api-review/SKILL.md   example skill (slide 18-19)
```

Request flow to read in code: request -> handler -> store / ProductClient -> JSON response.
Single `go.mod` for simplicity; each service is still built and run as its own binary.
