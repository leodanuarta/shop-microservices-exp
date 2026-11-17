# Code ownership challenges 

Batasan yang sengaja dibiarkan terbuka, supaya peserta bisa mengambil alih dari agent.

**Challenge A - status tidak dikenal.**
`GET /orders?status=xyz` sekarang mengembalikan `200 []`. Putuskan: tetap `200 []` atau `400`?
Dokumentasikan di `docs/SPEC.md`, implementasikan, tambahkan test. Lakukan hal yang sama untuk `category` di product-service.

**Challenge B - stok tidak berkurang.**
Setelah `POST /orders` sukses, stok di product-service tidak berubah, jadi order yang sama bisa dibuat
berulang sampai melebihi stok. Tulis test yang mereproduksi masalahnya, lalu rancang perbaikannya
(mis. endpoint reservasi stok), dan jelaskan trade-off-nya (race condition antar service).
Catatan: ini mengubah kontrak API - minta agent membuat rencana dulu dan review sebelum menerapkan.

**Challenge C - review diff.**
Minta agent "tambahkan pagination ke GET /products". Review diff-nya: tolak perubahan di luar scope
(refactor, dependency baru, file tidak terkait) dan jelaskan alasannya.

**Output:** demo singkat + hasil `go test ./...` + penjelasan bagian kode yang diubah.
