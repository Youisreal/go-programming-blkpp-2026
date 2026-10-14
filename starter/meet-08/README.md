# Northwind CRUD — Go + Supabase PostgreSQL

REST API CRUD produk dengan `net/http`, `database/sql`, dan raw SQL PostgreSQL. Tanpa framework atau ORM; satu dependency Go: driver `github.com/lib/pq`.

Koneksi langsung ke PostgreSQL Supabase, bukan REST API Supabase. Gunakan connection string database, bukan anon key atau service role API key. `database/sql` menyediakan connection pool; driver PostgreSQL tetap diperlukan.

`northwind.sql` berisi **subset data Northwind**: 8 kategori, 3 supplier, 8 produk. Ini bukan dump Northwind lengkap. File asli belum tersedia saat project dibuat. SQL sekarang memakai sintaks PostgreSQL; nama kolom dump lain perlu disesuaikan dengan project.

## Struktur kode

```text
main.go                         startup dan perakitan dependency
internal/config/config.go       environment, validasi URL database, SSL
internal/database/postgres.go   koneksi PostgreSQL, pool, ping, timeout
internal/product/product.go     input produk dan error domain
internal/product/repository.go  raw SQL dan pemetaan error PostgreSQL
internal/product/service.go     normalisasi, validasi, proses CRUD
internal/product/handler.go     routing HTTP, decode JSON, status/response
northwind.sql                   schema dan data sampel
```

Handler memanggil service, service memanggil repository, repository mengakses `*sql.DB`. Ubah aturan produk di service, query di repository, atau konfigurasi koneksi di config/database. Dependency dirakit langsung di `main.go`; tidak ada framework DI atau interface untuk satu implementasi repository. `go run .` tetap berlaku.

## Jalankan

Butuh Go 1.20+ dan project Supabase.

1. Buka dashboard Supabase, klik **Connect**, pilih connection string PostgreSQL. Untuk jaringan IPv4 pilih **Session pooler**, port `5432`. Direct connection `db.PROJECT_REF.supabase.co` membutuhkan IPv6 kecuali project memakai dukungan IPv4. Salin host, username, dan port persis dari dashboard; gunakan password database.
2. Buka **SQL Editor**, jalankan isi `northwind.sql` untuk membuat tabel dan data sampel. Tabel memakai RLS tanpa policy publik; koneksi backend contoh memakai role `postgres` yang dapat mengaksesnya. Jika memakai role lain, atur izin dan policy sesuai kebutuhan.
3. Salin konfigurasi contoh dan isi connection string lewat editor lokal. `.env` diabaikan Git. URL-encode karakter khusus password, misalnya `@` menjadi `%40`.

   ```sh
   cp .env.example .env
   ```

4. Muat konfigurasi lalu jalankan. `.env` tidak dimuat otomatis oleh Go:

```sh
set -a
source .env
set +a
go mod download
go run .
```

Server berjalan di `http://127.0.0.1:8080`. Startup memeriksa koneksi dan tabel produk. SSL wajib: default `sslmode=require`; `verify-ca`/`verify-full` juga didukung. Untuk verifikasi sertifikat/host gunakan `verify-full` dengan konfigurasi CA Supabase yang sesuai. URL/password tidak ditulis ke log startup.

Import dilakukan manual, bukan setiap startup. Import ulang menambah ID sampel yang belum ada, termasuk produk sampel yang sebelumnya dihapus; record yang sudah ada tidak ditimpa. Jangan import ulang jika ingin mempertahankan penghapusan produk sampel. Sequence ID disesuaikan agar POST tidak bentrok dengan ID seed.

Alternatif import dengan PostgreSQL CLI setelah `DATABASE_URL` dimuat:

```sh
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f northwind.sql
```

`DATABASE_URL` wajib tersedia. `ADDR` opsional untuk mengganti alamat server. File SQLite lama tidak dipakai dan tidak dimigrasikan; import SQL ke Supabase untuk memulai data sampel baru.

## Endpoint

| Method | Path | Fungsi |
| --- | --- | --- |
| GET | `/products` | Daftar produk |
| GET | `/products/{id}` | Detail produk |
| POST | `/products` | Tambah produk |
| PUT | `/products/{id}` | Ganti data produk |
| DELETE | `/products/{id}` | Hapus produk |
| GET | `/categories` | Daftar kategori untuk memilih `category_id` |
| GET | `/suppliers` | Daftar supplier untuk memilih `supplier_id` |

Daftar mengembalikan array JSON; detail/create/update/delete mengembalikan satu objek produk. POST mengembalikan 201, operasi berhasil lainnya 200. Error berupa `{"error":"pesan"}`: 400 untuk input salah, 404 untuk produk/endpoint tidak ditemukan, 405 untuk method salah, 409 untuk foreign key salah, 500 untuk kegagalan database.

## Contoh CRUD

1. Lihat data dan pilihan kategori/supplier:

   ```sh
   curl http://127.0.0.1:8080/products
   curl http://127.0.0.1:8080/products/1
   curl http://127.0.0.1:8080/categories
   curl http://127.0.0.1:8080/suppliers
   ```

2. Tambah produk:

   ```sh
   curl -i -X POST http://127.0.0.1:8080/products \
     -H 'Content-Type: application/json' \
     -d '{"product_name":"Teh Melati","supplier_id":1,"category_id":1,"quantity_per_unit":"10 boxes","unit_price":15.5,"units_in_stock":20,"units_on_order":0,"reorder_level":5,"discontinued":0}'
   ```

3. Gunakan `product_id` dari respons POST. Pada database baru, ID pertama adalah 9:

   ```sh
   curl -i -X PUT http://127.0.0.1:8080/products/9 \
     -H 'Content-Type: application/json' \
     -d '{"product_name":"Teh Melati Baru","supplier_id":1,"category_id":1,"quantity_per_unit":"10 boxes","unit_price":16,"units_in_stock":25,"units_on_order":0,"reorder_level":5,"discontinued":0}'

   curl -i -X DELETE http://127.0.0.1:8080/products/9
   ```

PUT mengganti seluruh data, bukan PATCH. `product_name`, `supplier_id`, dan `category_id` wajib valid. Field lain jika tidak dikirim memakai nilai nol/string kosong. Nama maksimal 40 karakter, kemasan maksimal 20, harga 0–9999999999.99, stok/order/reorder 0–32767, `discontinued` 0 atau 1. Database menyimpan harga sebagai `NUMERIC(12,2)`; harga dibulatkan ke dua digit desimal oleh PostgreSQL.

SQL memakai parameter `$1`, `$2`, dan seterusnya. Nilai request dikirim terpisah dari query lewat driver, sehingga kutip/newline tetap data. PostgreSQL menegakkan foreign key dan CHECK constraint.

API untuk latihan lokal, tanpa autentikasi. Default hanya mendengarkan loopback.

## Cek

```sh
go test -v ./...
go vet ./...
go build .
```

`go test` memeriksa config/environment, SSL, serta validasi service tanpa database. Untuk integration check CRUD, validasi, foreign key, SQL injection, dan import ulang, siapkan PostgreSQL test lalu jalankan:

```sh
TEST_DATABASE_URL='postgresql://USER:PASSWORD@localhost:5432/TEST_DB?sslmode=disable' go test -v ./...
```

`sslmode=disable` hanya untuk PostgreSQL test lokal; aplikasi menolaknya. Test membuat schema unik sementara dan menghapus schema tersebut setelah selesai. Gunakan database test dengan izin CREATE SCHEMA dan role superuser/BYPASSRLS agar dapat mengakses seed yang mengaktifkan RLS. Gunakan direct/session connection; transaction pooler tidak mendukung session `search_path` test. Tanpa `TEST_DATABASE_URL`, test CRUD dilewati.
