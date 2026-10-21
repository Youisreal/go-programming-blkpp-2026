# Product CRUD API — Gin Framework

CRUD REST API dengan Gin framework menggunakan arsitektur per-fitur (Handler → Service → Repository).

## 1. Tech Stack & Prasyarat

- Go 1.22+
- Gin Web Framework v1.12.0
- Storage: In-memory (thread-safe dengan `sync.RWMutex`)

## 2. Struktur Proyek

```text
meet-10/
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── README.md
└── internal/
    └── product/
        ├── model.go       # Struct data dan validasi input
        ├── repository.go  # Akses data (in-memory map)
        ├── service.go     # Logika bisnis
        ├── handler.go     # HTTP handler (Gin context)
        └── routes.go      # Registrasi endpoint
```

## 3. Setup & Menjalankan

### Download dependency (jika belum)
```bash
go mod tidy
```

### Jalankan server
```bash
go run main.go
```
Server berjalan di `http://localhost:8080`.

Port custom lewat environment variable:
```bash
PORT=3000 go run main.go
```

### Jalankan unit test
```bash
go test -v ./...
```

## 4. Panduan Penggunaan API

### Health Check
```bash
curl -i http://localhost:8080/health
```

### 1. Tambah Produk (Create)
```bash
curl -i -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Kopi Tubruk", "price": 15000}'
```
Response `201 Created`:
```json
{
  "data": {
    "id": 1,
    "name": "Kopi Tubruk",
    "price": 15000
  }
}
```

### 2. Ambil Semua Produk (Read All)
```bash
curl -i http://localhost:8080/api/products
```
Response `200 OK`:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Kopi Tubruk",
      "price": 15000
    }
  ]
}
```

### 3. Ambil Produk berdasarkan ID (Read One)
```bash
curl -i http://localhost:8080/api/products/1
```
Response `200 OK`:
```json
{
  "data": {
    "id": 1,
    "name": "Kopi Tubruk",
    "price": 15000
  }
}
```

### 4. Perbarui Produk (Update)
```bash
curl -i -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Kopi Susu", "price": 18000}'
```
Response `200 OK`:
```json
{
  "data": {
    "id": 1,
    "name": "Kopi Susu",
    "price": 18000
  }
}
```

### 5. Hapus Produk (Delete)
```bash
curl -i -X DELETE http://localhost:8080/api/products/1
```
Response `204 No Content`.

## 5. Format Error

Response error seragam:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Produk tidak ditemukan"
  }
}
```
Kode error:
- `400`: `INVALID_ID`, `INVALID_JSON`, `INVALID_INPUT`
- `404`: `NOT_FOUND`
- `500`: `INTERNAL_ERROR`
