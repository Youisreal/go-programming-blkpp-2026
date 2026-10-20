# Gin Auth + GORM

Project autentikasi dengan Gin, GORM, Neon PostgreSQL, bcrypt, dan JWT HS256.
Struktur per fitur mengikuti `../basic-orm`: Handler, Service, Repository.
Tabel `users` dibuat otomatis lewat `AutoMigrate` ketika server dimulai.

## Menjalankan

Go 1.25+ dan database PostgreSQL di Neon. Driver PostgreSQL memakai pgx;
compiler C tidak diperlukan.

```bash
cd gin-auth
go mod download
cp .env.example .env
```

Di Neon Dashboard, buka project lalu pilih **Connect**. Salin connection string
ke `DATABASE_URL` di `.env`. Pilih host `-pooler` untuk koneksi pooled atau host
langsung untuk koneksi direct. Contoh format, bukan kredensial aktif:

```dotenv
DATABASE_URL='postgresql://USER:PASSWORD@ep-PROJECT-pooler.REGION.aws.neon.tech/neondb?sslmode=require&channel_binding=require'
JWT_SECRET='hasil-openssl-rand-hex-32'
PORT=8080
```

Isi `JWT_SECRET` dengan hasil `openssl rand -hex 32`, minimal 32 byte.
Jangan commit `.env`. `.env` dibaca oleh shell lewat perintah berikut:

```bash
set -a
source .env
set +a
go run .
```

Default server: `http://localhost:8080`. `DATABASE_URL` wajib diisi dan harus memakai
TLS (`sslmode=require`, `verify-ca`, atau `verify-full`). Gunakan `verify-full`
untuk verifikasi sertifikat dan hostname. Password dengan karakter khusus dalam
URL harus di-encode; connection string dari Neon sudah memakai format URL.

Koneksi diperiksa saat startup, dengan timeout default 10 detik. Pool dibatasi
10 koneksi aktif dan 5 idle. `PreferSimpleProtocol` menghindari cache prepared
statement pada koneksi pooler. Aplikasi berhenti jika konfigurasi, koneksi,
atau migrasi gagal; tidak ada fallback ke database lain.

## Endpoint dan middleware

Semua request melewati logger, recovery, serta pembatas body 1 MiB.
Setiap endpoint memiliki hook `MockEndpoint` sendiri di `internal/middleware/mock.go`.

| Method | Endpoint | Header `X-Mock-Middleware` | Auth |
| --- | --- | --- | --- |
| GET | `/health` | `health` | Publik |
| POST | `/api/auth/register` | `auth.register` | Publik |
| POST | `/api/auth/login` | `auth.login` | Publik |
| GET | `/api/me` | `auth.me` | Bearer JWT |

Mock hanya menambahkan header untuk menunjukkan middleware berjalan.
Tambahkan logika sebelum/sesudah `c.Next()` atau kondisi berdasarkan `name`
untuk latihan aturan tiap endpoint. Middleware `RequireAuth` memvalidasi JWT
dan mengecek pengguna di database. Mock tidak melewati validasi auth.
Route auth dan middleware terdaftar di `internal/auth/routes.go`;
route health dan perakitan dependency berada di `setupEngine` pada `main.go`.

## Contoh request

```bash
curl -i http://localhost:8080/health

curl -i -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Budi","email":"budi@example.com","password":"password123"}'

curl -i -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"budi@example.com","password":"password123"}'
```

Login mengembalikan bentuk berikut; salin `access_token` ke request profil:

```json
{"data":{"access_token":"<JWT>","token_type":"Bearer","expires_in":3600}}
```

```bash
TOKEN='<access_token dari login>'
curl -i http://localhost:8080/api/me -H "Authorization: Bearer $TOKEN"
```

Register mengembalikan `201` dan profil pengguna tanpa password/hash.
Password minimal 8 karakter dan maksimal 72 byte; email dinormalisasi ke huruf kecil
dan dibatasi unique index. Login/profil berhasil mengembalikan `200`.
Input salah mendapat `400`, email duplikat `409`, kredensial atau token salah `401`,
dan kegagalan database `500`. Format error: `{"error":"pesan"}`.

JWT berlaku satu jam, memerlukan issuer `gin-auth`, dan hanya menerima HS256.
Untuk latihan ini belum ada refresh token atau logout server-side;
tambahkan revocation store saat membutuhkan pencabutan token sebelum kedaluwarsa.

## Verifikasi

```bash
go test ./...
go vet ./...
```

Test tanpa database mencakup validasi URL/TLS, hook endpoint publik,
input tidak valid, serta penolakan token salah/kedaluwarsa.

Untuk test integrasi, buat branch/database Neon khusus test, lalu isi
`TEST_DATABASE_URL` dalam `.env` dengan connection string database tersebut.
Muat ulang `.env` dengan perintah `set -a`, `source .env`, `set +a` seperti di atas:

```bash
go test -v ./...
```

Test integrasi menjalankan migrasi tabel `users`, register/login/profil,
hash password, email duplikat, serta pembacaan lewat koneksi PostgreSQL terpisah.
Pengguna test memakai email unik dan dihapus setelah test. Test integrasi dilewati
jika `TEST_DATABASE_URL` kosong; keberhasilan test tanpa URL bukan bukti koneksi Neon.

## File

```text
gin-auth/
├── main.go              # Perakitan dependency, migrasi, server
├── main_test.go         # Test route dan integrasi PostgreSQL
├── internal/
│   ├── auth/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── handler.go
│   │   ├── routes.go
│   │   └── middleware.go
│   ├── database/
│   │   ├── database.go      # Koneksi Neon PostgreSQL
│   │   └── database_test.go
│   └── middleware/
│       └── mock.go         # Hook mock tiap endpoint
├── .env.example         # DATABASE_URL, JWT_SECRET, PORT, TEST_DATABASE_URL
├── .gitignore
├── go.mod
└── go.sum
```
