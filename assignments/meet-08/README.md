# Northwind CRUD — Go + Supabase PostgreSQL

REST API CRUD produk dengan `net/http`, `database/sql`, dan raw SQL PostgreSQL. Tanpa framework atau ORM; satu dependency Go: driver `github.com/lib/pq`.

Koneksi langsung ke PostgreSQL Supabase, bukan REST API Supabase. Gunakan connection string database, bukan anon key atau service role API key. `database/sql` menyediakan connection pool; driver PostgreSQL tetap diperlukan.

`northwind.sql` berisi **subset data Northwind** untuk PostgreSQL: 8 kategori, 3 supplier, 8 produk. Untuk memakai dump Northwind lain, sesuaikan schema dan nama kolom dengan query repository.

## Struktur kode

```text
main.go                         startup dan perakitan dependency
internal/config/config.go       environment, validasi URL database, SSL
internal/config/env.go          pembaca .env tanpa shell/dependency tambahan
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
3. Salin `.env.example` menjadi `.env` lewat editor/file manager, lalu isi connection string. Pastikan nama file `.env` dan encoding UTF-8. `.env` diabaikan Git. URL-encode karakter khusus password, misalnya `@` menjadi `%40`.
4. Jalankan dari folder project. Berlaku di Windows, macOS, dan Linux:

```sh
go run .
```

Server tersedia di `http://127.0.0.1:8080`. Hentikan dengan `Ctrl+C`.

## Konfigurasi

Contoh `.env` untuk Session pooler. Ganti placeholder dengan nilai dari dashboard Supabase:

```dotenv
DATABASE_URL='postgresql://postgres.PROJECT_REF:PASSWORD@POOLER_HOST:5432/postgres?sslmode=require'
ADDR='127.0.0.1:8080'
```

| Variabel | Kegunaan | Default |
| --- | --- | --- |
| `DATABASE_URL` | Connection string PostgreSQL | Wajib diisi |
| `ADDR` | Alamat HTTP server | `127.0.0.1:8080` |

Aplikasi membaca `.env` otomatis dari **folder kerja saat ini**, tanpa menjalankan shell. Environment OS lebih diutamakan, termasuk variabel yang sengaja diisi kosong. `.env` opsional jika konfigurasi tersedia dari OS; nilai file tidak diekspor ke environment proses. Connection string wajib berupa URL `postgresql://USER:PASSWORD@HOST:PORT/postgres?sslmode=require`, bukan anon key Supabase.

Format file: `KEY=value`, kutip tunggal/ganda literal, baris kosong, komentar pada baris tersendiri, dan prefix `export ` opsional. Mendukung UTF-8, BOM UTF-8, serta LF/CRLF. Simpan file sebagai UTF-8, bukan UTF-16. Tidak mendukung multiline, escape/interpolasi, atau komentar di akhir nilai. `$()` dan `${...}` dibaca sebagai teks; tidak dieksekusi. Input rusak menghasilkan nomor baris tanpa menampilkan nilai/password.

Startup memeriksa koneksi dan tabel produk. SSL wajib: default `sslmode=require`; `verify-ca`/`verify-full` juga didukung. Untuk verifikasi sertifikat/host gunakan `verify-full` dengan konfigurasi CA Supabase yang sesuai. Error koneksi menampilkan penyebab dengan connection string dan password disamarkan.

## Data sampel

Import dilakukan manual, bukan setiap startup. Import ulang menambah ID sampel yang belum ada, termasuk produk sampel yang sebelumnya dihapus; record yang sudah ada tidak ditimpa. Jangan import ulang jika ingin mempertahankan penghapusan produk sampel. Sequence ID disesuaikan agar POST tidak bentrok dengan ID seed.

Import SQL melalui SQL Editor Supabase. Aplikasi memakai PostgreSQL; file SQLite lama tidak dimigrasikan.

## Routing API

Routing memakai `net/http` bawaan Go. Di [main.go](main.go), handler produk dipasang langsung sebagai `Handler` milik HTTP server:

```go
service := product.Service{Repository: repository}
handler := product.Handler{Service: service}
server := &http.Server{Addr: cfg.Addr, Handler: handler}
```

`product.Handler` memenuhi interface `http.Handler` karena memiliki method `ServeHTTP(w http.ResponseWriter, r *http.Request)`. Server otomatis memanggil method ini untuk setiap request. Pemetaan route berada di [internal/product/handler.go](internal/product/handler.go), bukan di `main.go`.

Alur penanganan request:

1. `ServeHTTP` membaca `r.URL.Path`. Switch pertama mengenali `/products`, prefix `/products/`, `/categories`, dan `/suppliers`.
2. Untuk `/products/{id}`, bagian setelah `/products/` diambil dengan `strings.TrimPrefix`, lalu diubah menjadi integer positif melalui `strconv.Atoi`.
3. Untuk produk, switch kedua membaca `r.Method` dan memilih method service. Kategori/supplier ditangani langsung pada switch pertama dan hanya menerima GET.
4. Service menjalankan validasi/proses CRUD, repository menjalankan raw SQL, lalu handler menulis status HTTP dan JSON melalui `respond`.

Contoh pemetaan GET produk di `ServeHTTP`:

```go
case r.Method == http.MethodGet && collection:
    rows, err = h.Service.List(r.Context())
case r.Method == http.MethodGet:
    rows, err = h.Service.Get(r.Context(), id)
```

`collection` bernilai `true` ketika path tepat `/products`. Pada `GET /products/1`, nilainya `false` dan `id` bernilai `1`, sehingga handler memanggil `Service.Get`.

### Pemetaan endpoint

| Method | Path | Method service | Fungsi |
| --- | --- | --- | --- |
| GET | `/products` | `List(ctx)` | Daftar produk |
| GET | `/products/{id}` | `Get(ctx, id)` | Detail produk |
| POST | `/products` | `Create(ctx, input)` | Tambah produk |
| PUT | `/products/{id}` | `Update(ctx, id, input)` | Ganti data produk |
| DELETE | `/products/{id}` | `Delete(ctx, id)` | Hapus produk |
| GET | `/categories` | `Categories(ctx)` | Daftar kategori untuk memilih `category_id` |
| GET | `/suppliers` | `Suppliers(ctx)` | Daftar supplier untuk memilih `supplier_id` |

Daftar mengembalikan array JSON; detail/create/update/delete mengembalikan satu objek produk. POST mengembalikan 201, operasi berhasil lainnya 200. Error berupa `{"error":"pesan"}`: 400 untuk input salah, 404 untuk produk/endpoint tidak ditemukan, 405 untuk method salah, 409 untuk foreign key salah, 500 untuk kegagalan database.

Path yang tidak dikenali menghasilkan 404. ID produk yang kosong/tidak valid menghasilkan 400; misalnya `/products/`, `/products/abc`, atau `/products/1/`. Method yang tidak didukung pada path valid menghasilkan 405 disertai header `Allow`. Contohnya, `PATCH /products/1` mengembalikan `Allow: GET, PUT, DELETE`.

### Menambah route

1. Tambahkan query yang diperlukan di [repository.go](internal/product/repository.go), lalu proses/validasi di [service.go](internal/product/service.go).
2. Tambahkan pengenalan path dan pemetaan method di `ServeHTTP`. Untuk path spesifik seperti `/products/search`, tempatkan pemeriksaan **sebelum** prefix `/products/` agar `search` tidak dianggap sebagai ID.
3. Gunakan `respond` untuk hasil service. Untuk method yang ditolak, isi header `Allow` dan kembalikan 405. Untuk input JSON, pertahankan batas body dan validasi decoder seperti endpoint POST/PUT.
4. Tambahkan kasus integration test di [main_test.go](main_test.go) dan perbarui tabel endpoint di README.

## Contoh CRUD

### macOS/Linux — curl

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

### Windows — PowerShell

Gunakan `Invoke-RestMethod` agar contoh tidak bergantung pada alias `curl` atau aturan kutip command shell:

```powershell
$baseUrl = 'http://127.0.0.1:8080'
Invoke-RestMethod -Uri "$baseUrl/products"
Invoke-RestMethod -Uri "$baseUrl/categories"
Invoke-RestMethod -Uri "$baseUrl/suppliers"

$product = @{
    product_name = 'Teh Melati'
    supplier_id = 1
    category_id = 1
    quantity_per_unit = '10 boxes'
    unit_price = 15.5
    units_in_stock = 20
    units_on_order = 0
    reorder_level = 5
    discontinued = 0
}
$body = $product | ConvertTo-Json
$created = Invoke-RestMethod -Method Post -Uri "$baseUrl/products" -ContentType 'application/json' -Body $body
$id = $created.product_id
Invoke-RestMethod -Uri "$baseUrl/products/$id"

$product.product_name = 'Teh Melati Baru'
$product.unit_price = 16
$body = $product | ConvertTo-Json
Invoke-RestMethod -Method Put -Uri "$baseUrl/products/$id" -ContentType 'application/json' -Body $body
Invoke-RestMethod -Method Delete -Uri "$baseUrl/products/$id"
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

`go test` memeriksa pemuatan `.env`/environment, BOM/CRLF, kutip literal, perlindungan password pada error, SSL, serta validasi service tanpa database. Untuk integration check CRUD, validasi, foreign key, SQL injection, dan import ulang, siapkan PostgreSQL test.

macOS/Linux:

```sh
TEST_DATABASE_URL='postgresql://USER:PASSWORD@localhost:5432/TEST_DB?sslmode=disable' go test -v ./...
```

Windows PowerShell:

```powershell
$env:TEST_DATABASE_URL = 'postgresql://USER:PASSWORD@localhost:5432/TEST_DB?sslmode=disable'
go test -v ./...
Remove-Item Env:TEST_DATABASE_URL
```

`sslmode=disable` hanya untuk PostgreSQL test lokal; aplikasi menolaknya. Test membuat schema unik sementara dan menghapus schema tersebut setelah selesai. Gunakan database test dengan izin CREATE SCHEMA dan role superuser/BYPASSRLS agar dapat mengakses seed yang mengaktifkan RLS. Gunakan direct/session connection; transaction pooler tidak mendukung session `search_path` test. Tanpa `TEST_DATABASE_URL`, test CRUD dilewati.

## Masalah umum

| Pesan error | Perbaikan |
| --- | --- |
| `DATABASE_URL belum diisi` | Isi `.env` di folder kerja, atau atur environment OS. Periksa apakah variabel OS sudah ada tetapi kosong. |
| `DATABASE_URL wajib berupa connection string PostgreSQL lengkap` | Salin URI PostgreSQL lengkap dari dashboard, termasuk user, host, port, dan nama database. |
| `.env baris ...` | Perbaiki format `KEY=value`, pasangan kutip, dan encoding UTF-8 pada baris tersebut. |
| `lookup ...: no such host` | Periksa hostname. Pisahkan port dengan `:`, misalnya `.supabase.co:5432`. |
| `network is unreachable` atau timeout | Periksa jaringan/firewall. Untuk jaringan IPv4 gunakan Session pooler dari dashboard Supabase. |
| `password authentication failed` | Periksa username/password database dan URL-encode karakter khusus password. |
| `database belum siap; import northwind.sql dahulu` | Import SQL ke database yang dituju dan pastikan role koneksi memiliki izin mengakses tabel. |
