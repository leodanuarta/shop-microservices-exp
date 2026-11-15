# Spec (natural language -> testable behavior)

Pola dari slide: **Konteks + tujuan + input/output + batasan + acceptance criteria + verifikasi.**

## Context
Dua microservice Go kecil untuk e-commerce. Stdlib saja, data in-memory, tanpa database.

## Scope (hanya 5 endpoint bisnis + 2 health)
| Service | Endpoint | Fungsi |
|---|---|---|
| product-service :8081 | `GET /health` | health check |
| | `GET /products?category=` | list produk, filter kategori opsional |
| | `GET /products/{id}` | detail produk (dipakai order-service) |
| order-service :8082 | `GET /health` | health check |
| | `GET /orders?status=` | list order, filter status opsional |
| | `POST /orders` | buat order; validasi ke product-service |

Non-goals: auth, database, pembayaran, pagination, deployment, arsitektur production-grade.

## Acceptance criteria
**GET /products**
- Return JSON, `Content-Type: application/json`.
- Tanpa `category` -> semua produk. Dengan `category` -> hanya yang cocok (case-insensitive).
- Tidak ada yang cocok -> `200` dan `[]` (bukan `null`).

**GET /products/{id}** -> `200` bila ada, `404` bila tidak.

**POST /orders** body `{"items":[{"product_id":"p-1","quantity":2}]}`
- `201` + order (`status: pending`, `total` = sum(price x qty), IDR bulat).
- `400` JSON rusak / items kosong / quantity <= 0 / product_id kosong.
- `422` produk tidak dikenal. `409` stok tidak cukup. `502` product-service tidak bisa dihubungi.

**GET /orders?status=** -> sama seperti filter produk; status tak dikenal saat ini -> `200 []` (lihat Challenge A).

## Constraints
- Jangan tambah dependency. Jangan ubah file yang tidak relevan.
- Laporkan file yang diubah dan hasil test yang benar-benar dijalankan.

## Verification
`go vet ./... && go test ./...`, lalu jalankan kedua service dan coba dengan curl (lihat README).
