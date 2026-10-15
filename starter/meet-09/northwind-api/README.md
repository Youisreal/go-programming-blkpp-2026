# Northwind REST API — Workshop Go

## 1. Overview

REST API sederhana untuk database Northwind PostgreSQL schema `public` yang sudah ada di Supabase. Semua module memakai Handler → Service → Repository agar pola workshop konsisten. Service Product dan Category meneruskan operasi CRUD; Service Order menangani validasi bisnis, perhitungan, dan transaksi stok. Tidak ada migration otomatis atau pembuatan ulang tabel.

## 2. Tujuan pembelajaran

- Membuat routing HTTP dengan method dan path parameter Go 1.22.
- Membaca JSON, memvalidasi input, dan menghasilkan status HTTP konsisten.
- Menulis query SQL parameterized dan menangani error PostgreSQL.
- Memahami foreign key, transaksi, row locking, dan rollback.
- Mengembangkan module baru tanpa framework atau generic repository.

## 3. Tech stack

Go minimal 1.22, `net/http`, `database/sql`, PostgreSQL, `lib/pq` v1.10.9, dan `godotenv` v1.5.1. Test memakai `testing` dan `httptest`. Hanya dua dependency eksternal.

## 4. Struktur folder

```text
northwind-api/
├── main.go
├── main_test.go
├── internal/
│   ├── database/postgres.go
│   ├── httpjson/
│   │   ├── httpjson.go
│   │   └── httpjson_test.go
│   ├── product/
│   │   ├── model.go
│   │   ├── model_test.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── handler.go
│   │   └── routes.go
│   ├── category/
│   │   ├── model.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── handler.go
│   │   └── routes.go
│   └── order/
│       ├── model.go
│       ├── repository.go
│       ├── service.go
│       ├── service_test.go
│       ├── handler.go
│       └── routes.go
├── examples/requests.http
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

`internal/httpjson` merupakan helper kecil bersama untuk envelope JSON, decode, validasi ID, dan pemetaan error database. Tidak ada global DB atau interface tambahan.

## 5. Tanggung jawab layer

| File | Tanggung jawab |
|---|---|
| `main.go` | Load environment, koneksi, explicit dependency wiring, HTTP server |
| `model.go` | Struct, JSON tags, validasi nilai, perhitungan murni |
| `routes.go` | Daftarkan method dan path ke `http.ServeMux` |
| `handler.go` | Baca request/path, validasi HTTP, panggil service, tulis response |
| `repository.go` | SQL, parameter, Scan, RowsAffected, Close, dan error database |
| `service.go` | Menjalankan use case; Product/Category meneruskan CRUD ke repository, Order mengatur aturan bisnis dan lifecycle transaksi |

Handler meneruskan `r.Context()` sampai query. Semua handler memakai service, lalu service memakai repository. Semua dependency berupa concrete struct dan diberikan melalui constructor.

### Pembagian package dan dependency

Project memakai struktur berdasarkan fitur: kode Product berada di `internal/product`, Category di `internal/category`, dan Order di `internal/order`. Setiap module memiliki model, routes, handler, service, dan repository sendiri. Service Product/Category sengaja tipis agar siswa mempelajari satu alur yang sama; Service Order menunjukkan kebutuhan aturan bisnis lintas tabel.

| Package | Dipakai oleh | Peran |
|---|---|---|
| `internal/database` | `main.go` | Membuka koneksi PostgreSQL, mengatur pool, dan memeriksa koneksi saat startup |
| `internal/product` | `main.go` | CRUD product melalui handler, service, dan repository |
| `internal/category` | `main.go` | CRUD category melalui handler, service, dan repository |
| `internal/order` | `main.go` | Read order dan create order melalui handler, service, dan repository |
| `internal/httpjson` | Handler dan routes setiap module, serta `main.go` | Format JSON, decode request, validasi path ID, pemetaan error database, dan fallback method HTTP |

Order mengakses tabel `public.products` lewat repository Order sendiri karena pemeriksaan harga/stok harus memakai transaksi yang sama dengan penyimpanan order. Order tidak memanggil handler atau repository Product. Module Product dan Category juga tidak saling memanggil; hubungan datanya dijaga foreign key PostgreSQL.

### Lokasi routes dan helper HTTP

Routes tetap berada di masing-masing module agar peserta dapat menelusuri satu fitur dalam satu folder. `main.go` memanggil `product.RegisterRoutes`, `category.RegisterRoutes`, dan `order.RegisterRoutes`; definisi method/path berada di `routes.go` masing-masing. Route `/health` dan fallback endpoint tidak ditemukan didaftarkan di `main.go`.

`httpjson` tetap berada di `internal` karena digunakan bersama dalam aplikasi ini. Go membatasi import package tersebut ke kode di pohon direktori induk `internal`, yaitu project `northwind-api`. Belum diperlukan package publik `pkg/` atau Go module terpisah.

| Helper | Kegunaan di batas HTTP |
|---|---|
| `JSON`, `Data`, `Error` | Menulis JSON dan envelope sukses/error secara konsisten |
| `ID` | Membaca `{id}` dan memvalidasi rentang ID |
| `Decode` | Membatasi ukuran body, decode satu nilai JSON, dan menolak field tidak dikenal |
| `DatabaseError` | Mengubah error akses data menjadi status HTTP tanpa mengirim detail SQL |
| `Fallbacks` | Menghasilkan response JSON 405 dan header `Allow` |

Repository mengembalikan data/error ke service tanpa mengetahui HTTP. Service meneruskan hasil ke handler; handler memakai `httpjson` untuk menerjemahkannya menjadi response. Helper ini tidak menjalankan query atau aturan bisnis.

### Pemakaian model dan pembagian validasi

`model.go` mendefinisikan Go struct sebagai wadah data, bukan layer yang dipanggil setelah repository. Karena file dalam satu direktori memakai package yang sama, handler, service, dan repository dapat memakai model tanpa import tambahan.

| Model | Dipakai untuk |
|---|---|
| `product.Product` | Menampung request create/update, menjadi parameter service/repository, menerima hasil `Scan`, dan membentuk response product |
| `category.Category` | Menampung request create/update, menjadi parameter service/repository, menerima hasil `Scan`, dan membentuk response category |
| `order.CreateRequest`, `order.ItemRequest` | Menampung JSON create order dan input aturan bisnis di service |
| `order.Order`, `order.Item` | Menampung hasil query atau hasil transaksi, termasuk harga, subtotal, dan total, untuk response |

```text
Request:  JSON → Decode di handler → struct → Service → Repository → SQL
Response: SQL → Scan di repository → struct → Service → Handler → JSON
```

`httpjson.Decode` membaca JSON request ke struct; `httpjson.Data` menulis nilai Go sebagai JSON response. Pemindahan hasil SQL ke struct dilakukan oleh `Scan` di repository. `httpjson` tidak mengubah hasil SQL secara langsung menjadi JSON.

Validasi struktur JSON dan path ID berada di handler melalui helper HTTP. Handler Product/Category memanggil `Valid()` pada model sebelum meneruskan create/update ke service. Service CRUD belum menambah validasi sendiri. Service Order memanggil `Validate()` untuk input order lalu memeriksa customer, produk, harga, dan stok dalam transaksi. Constraint PostgreSQL tetap menjaga PK, FK, dan aturan database.

### Wiring saat startup

`run()` di `main.go` memuat environment, memanggil `database.ConnectDatabase`, lalu memberikan satu `*sql.DB` ke `newMux(db)`. Fungsi `newMux` membuat repository, service setiap module, handler, dan router secara eksplisit:

```text
Product:  db → NewRepository → NewService → NewHandler → RegisterRoutes
Category: db → NewRepository → NewService → NewHandler → RegisterRoutes
Order:    db → NewRepository → NewService → NewHandler → RegisterRoutes
```

Panah di atas menunjukkan urutan pembuatan dependency saat startup. Semua repository berbagi pool `*sql.DB` yang sama. Service Order memiliki lifecycle `sql.Tx` per create request; koneksi pool dimiliki aplikasi dan ditutup ketika `run()` selesai. Setelah wiring selesai, `http.Server` menerima request melalui `http.ServeMux`.

## 6. Diagram arsitektur dan alur request

```mermaid
flowchart TD
    Client["Client: REST Client / curl"] --> Mux["net/http.ServeMux"]

    subgraph Product["internal/product"]
        PR["routes.go"] --> PH["handler.go"]
        PH --> PS["service.go"]
        PS --> PRepo["repository.go"]
    end

    subgraph Category["internal/category"]
        CR["routes.go"] --> CH["handler.go"]
        CH --> CS["service.go"]
        CS --> CRepo["repository.go"]
    end

    subgraph Order["internal/order"]
        OR["routes.go"] --> OH["handler.go"]
        OH --> OS["service.go"]
        OS --> ORepo["repository.go"]
    end

    Mux --> PR
    Mux --> CR
    Mux --> OR
    Mux --> Health["GET /health: tanpa query DB"]

    PH -.-> JSON["internal/httpjson: helper HTTP bersama"]
    CH -.-> JSON
    OH -.-> JSON

    PRepo --> SQL["database/sql + lib/pq"]
    CRepo --> SQL
    ORepo --> SQL
    SQL --> DB[("Supabase PostgreSQL: schema public")]
```

Panah penuh menunjukkan alur request/akses data. Panah putus-putus menunjukkan pemakaian helper HTTP oleh handler. `model.go` di setiap module mendefinisikan bentuk data yang dipakai layer tersebut; model bukan tahap tambahan dalam alur request. Diagram menampilkan jalur bisnis utama; routes juga memakai helper untuk fallback method.

Versi teks untuk pembaca tanpa dukungan Mermaid:

```text
HTTP Request → Routes → Handler → Service → Repository → PostgreSQL
```

Contoh `GET /api/products/1`:

1. `ServeMux` mencocokkan method/path yang didaftarkan `product/routes.go` dan mengisi path value `id`.
2. `product.Handler.GetByID` membaca ID melalui `httpjson.ID`.
3. Handler memanggil `product.Service.GetByID` dengan context request.
4. Service meneruskan pemanggilan ke `product.Repository.GetByID`.
5. Repository menjalankan query parameterized pada `public.products` lalu melakukan `Scan`.
6. Hasil kembali melalui service ke handler. Handler mengirim `httpjson.Data` dengan status 200, atau memetakan error melalui `httpjson.DatabaseError`.

Semua module memakai urutan layer yang sama. Perbedaannya terletak pada isi service: CRUD sederhana diteruskan, sedangkan Order menjalankan aturan bisnis dan transaksi.

## 7. Service Layer untuk pola workshop yang konsisten

Di project ini, semua module memiliki service agar siswa mengikuti pola yang sama. Product dan Category belum memiliki aturan bisnis tambahan, sehingga service cukup meneruskan CRUD ke repository. Pada aplikasi kecil lain, layer penerusan ini boleh dilewati; di workshop ini konsistensi alur menjadi alasan mempertahankannya. Jangan menambahkan aturan bisnis buatan hanya untuk mengisi service.

Order menunjukkan peran service ketika satu request melibatkan beberapa tabel dan aturan stok:

1. Validasi ID, customer ID, items, quantity, dan discount; tolak product ID berulang.
2. Mulai satu `sql.Tx`; cek customer dengan `FOR KEY SHARE`.
3. Urutkan product ID lalu baca harga dan stok memakai `SELECT ... FOR UPDATE` di transaksi yang sama. Urutan lock konsisten mengurangi risiko deadlock antar-order.
4. Tolak produk hilang, discontinued, harga NULL/tidak valid, atau stok kurang. Stok NULL dianggap tidak tersedia.
5. Hitung `unit_price * quantity * (1 - discount)` dan jumlahkan total. Harga tidak diterima dari request.
6. Insert order, insert setiap detail, kurangi stok dengan conditional update.
7. Commit setelah semua sukses. `defer tx.Rollback()` membatalkan seluruh perubahan bila ada error; setelah commit pemanggilan rollback tidak mengubah hasil.

Lock produk bertahan sampai transaksi berakhir sehingga dua order tidak menjual stok yang sama. FK tetap diperiksa PostgreSQL. Total tidak disimpan sebagai kolom baru; GET menghitung dari detail. `GET /api/orders` menampilkan ringkasan, GET per ID menambahkan items terurut. Response create juga mengurutkan items berdasarkan product ID.

Northwind memakai `real` untuk harga dan diskon. `float64` dipakai **hanya untuk pembelajaran**, bukan representasi ideal keuangan production. Diskon create dikuantisasi ke `float32` agar sesuai kolom `real`; query Order membaca harga/diskon sebagai `double precision` agar hasil perhitungan konsisten dengan nilai tersimpan. Angka dapat memiliki galat pecahan biner. Untuk keuangan production, rancang tipe `numeric` atau satuan uang terkecil serta aturan pembulatan eksplisit.

## 8. Setup environment

Dari workspace:

```bash
cd northwind-api
cp -n .env.example .env
```

Opsi `-n` mempertahankan `.env` yang sudah ada. Edit `.env` secara lokal dan ganti placeholder host/password dengan konfigurasi Supabase:

```dotenv
PORT=8080
DATABASE_URL=postgresql://postgres:PASSWORD@HOST:5432/postgres?sslmode=require
```

`.env` diabaikan Git. Jangan masukkan kredensial ke source, request examples, atau README. Environment OS memiliki prioritas dibanding `.env`. `.env` tidak wajib jika variabel sudah diekspor. `DATABASE_URL` wajib saat menjalankan server; tidak diperlukan saat build/test. Gunakan Go 1.22 atau lebih baru.

## 9. Koneksi Supabase dan asumsi schema

Salin connection string PostgreSQL dari panel **Connect** project Supabase. Pilih direct connection atau **Session Pooler**, lalu gunakan host, user, password, database, dan port persis dari panel. User pooler biasanya mengandung project reference. Direct host dapat memerlukan IPv6; Session Pooler membantu jaringan IPv4. Gunakan `sslmode=require` dan URL-encode karakter khusus password seperti `@`, `:`, `/`, dan `#`. API ini memakai PostgreSQL langsung, bukan Supabase REST API/SDK atau anon key.

Schema live dan integration database **belum diverifikasi**. Keberadaan `.env` tidak berarti schema atau transaksi sudah diuji; unit test berjalan tanpa koneksi database. Asumsi Northwind PostgreSQL umum berikut dipakai:

| Tabel | Kolom yang digunakan dan asumsi |
|---|---|
| `public.products` | `product_id smallint PRIMARY KEY` tanpa default, `product_name varchar(40) NOT NULL`, `supplier_id smallint NULL` FK suppliers, `category_id smallint NULL` FK categories, `unit_price real NULL`, `units_in_stock smallint NULL`, `discontinued integer NOT NULL` bernilai 0/1 |
| `public.categories` | `category_id smallint PRIMARY KEY` tanpa default, `category_name varchar(15) NOT NULL`, `description text NULL` |
| `public.customers` | `customer_id varchar(5)/char(5) PRIMARY KEY`, contoh dataset `ALFKI` |
| `public.orders` | `order_id smallint PRIMARY KEY` tanpa default, `customer_id` nullable FK customers, `order_date date` nullable; API selalu mengisi customer dan tanggal create |
| `public.order_details` | `(order_id, product_id)` composite primary key dan FK, `unit_price real NOT NULL`, `quantity smallint NOT NULL`, `discount real NOT NULL` |

Kolom lain seperti `quantity_per_unit`, `units_on_order`, `reorder_level`, `picture`, employee, freight, dan shipping diasumsikan nullable atau memiliki default. API tidak mengelolanya. Tidak ada asumsi ID auto-increment: create selalu meminta ID positif 1–32767. Batas positif merupakan kebijakan API; PostgreSQL smallint sendiri juga menerima nilai negatif.

Product menerima supplier/category ID nullable; foreign key yang salah menghasilkan 409. Product GET memakai COALESCE untuk harga/stok NULL menjadi 0; PUT menyimpan harga/stok baru sebagai nilai non-NULL. Category description menggunakan pointer agar NULL tetap JSON `null`. `discontinued` wajib bertipe angka 0/1 pada varian ini; schema boolean memerlukan penyesuaian model/query. PUT mengganti seluruh field yang dikelola API; ID body yang tidak diisi atau bernilai 0 memakai ID path, sedangkan nilai selain 0 harus sama dengan path. Field angka yang tidak diberikan memakai nilai nol Go; `discontinued` tidak diberikan berarti 0. Field tabel yang tidak dikelola API tidak diubah saat PUT.

Sebelum memakai database berbeda, jalankan query read-only berikut lewat SQL Editor Supabase atau `psql` untuk memeriksa tipe, nullability, default, PK, dan FK:

```sql
SELECT table_name, column_name, data_type, udt_name,
       is_nullable, column_default, character_maximum_length
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name IN ('products', 'categories', 'customers', 'orders', 'order_details')
ORDER BY table_name, ordinal_position;

SELECT c.conrelid::regclass AS table_name, c.conname,
       pg_get_constraintdef(c.oid) AS definition
FROM pg_constraint c
JOIN pg_namespace n ON n.oid = c.connamespace
WHERE n.nspname = 'public'
  AND c.conrelid IN (
    'public.products'::regclass, 'public.categories'::regclass,
    'public.customers'::regclass, 'public.orders'::regclass,
    'public.order_details'::regclass)
ORDER BY table_name, c.conname;
```

Sesuaikan kode jika hasil berbeda, terutama `discontinued`, kolom NOT NULL tanpa default, dan tindakan FK `ON DELETE`. Penanganan 409 mengasumsikan relasi Northwind standar yang membatasi penghapusan data yang masih dipakai; schema dengan `ON DELETE CASCADE` memiliki perilaku berbeda.

## 10. Menjalankan aplikasi

Jalankan dari direktori `northwind-api` agar `.env` ditemukan:

```bash
go mod tidy
go run .
```

Startup menjalankan `PingContext` dengan timeout 10 detik. Pool maksimal 10 koneksi aktif, 5 idle, lifetime 30 menit, idle time 5 menit. Server default port 8080 dengan timeout HTTP. Server baru berjalan setelah koneksi berhasil. `GET /health` mengecek HTTP server tanpa query database per request; startup tetap memerlukan database valid. Health dapat diuji tanpa DB lewat unit test.

```bash
curl -i http://localhost:8080/health
```

## 11. Daftar endpoint

| Method | Endpoint | Status sukses |
|---|---|---|
| GET | `/health` | 200 |
| GET | `/api/products` | 200 |
| GET | `/api/products/{id}` | 200 |
| POST | `/api/products` | 201 |
| PUT | `/api/products/{id}` | 200 |
| DELETE | `/api/products/{id}` | 204 |
| GET | `/api/categories` | 200 |
| GET | `/api/categories/{id}` | 200 |
| POST | `/api/categories` | 201 |
| PUT | `/api/categories/{id}` | 200 |
| DELETE | `/api/categories/{id}` | 204 |
| GET | `/api/orders` | 200 |
| GET | `/api/orders/{id}` | 200 |
| POST | `/api/orders` | 201 |

ServeMux Go juga mendukung HEAD untuk route GET. Method lain menghasilkan 405 JSON dengan header Allow. Endpoint tidak ditemukan menghasilkan 404 JSON. Daftar diurutkan ID dan array kosong berbentuk `[]`.

## 12. Contoh request dan response

### 12.1. Aturan request bersama

Base URL lokal: `http://localhost:8080`. Contoh ID `30000` dipakai untuk data latihan; periksa dahulu agar tidak bentrok dengan data yang sudah ada.

- Gunakan `Content-Type: application/json` pada POST dan PUT.
- Kirim object JSON langsung sebagai body, tanpa pembungkus `data`. Pembungkus `data` hanya digunakan pada response sukses.
- Body maksimal 1 MiB dan hanya boleh berisi satu nilai JSON. Field yang tidak dikenal, tipe JSON yang tidak sesuai, dan JSON rusak menghasilkan 400 `INVALID_JSON`.
- Semua path `{id}` harus bilangan bulat positif 1–32767. Nilai di luar rentang/bukan angka menghasilkan 400 `INVALID_ID`.
- GET dan DELETE tidak memerlukan body. API belum menyediakan query parameter untuk filter, pencarian, atau pagination.
- Field bertipe integer ditulis sebagai angka JSON, misalnya `30000`, bukan string `"30000"`. `discontinued` memakai angka 0/1, bukan boolean.
- PUT mengganti seluruh field yang dikelola API. Field yang tidak dikirim tidak mempertahankan nilai lama: angka menjadi 0 dan pointer menjadi `null`. Kirim nilai lengkap yang ingin disimpan.

### 12.2. Ringkasan payload setiap endpoint

| Method dan endpoint | Path parameter | Body request | Sukses |
|---|---|---|---|
| `GET /health` | Tidak ada | Tidak ada | 200, status HTTP server |
| `GET /api/products` | Tidak ada | Tidak ada | 200, array product |
| `GET /api/products/{id}` | ID product | Tidak ada | 200, satu product |
| `POST /api/products` | Tidak ada | Object Product Create, lihat 12.3 | 201, product yang dibuat |
| `PUT /api/products/{id}` | ID product yang diubah | Object Product Update, lihat 12.3 | 200, product yang diperbarui |
| `DELETE /api/products/{id}` | ID product yang dihapus | Tidak ada | 204, tanpa body response |
| `GET /api/categories` | Tidak ada | Tidak ada | 200, array category |
| `GET /api/categories/{id}` | ID category | Tidak ada | 200, satu category |
| `POST /api/categories` | Tidak ada | Object Category Create, lihat 12.4 | 201, category yang dibuat |
| `PUT /api/categories/{id}` | ID category yang diubah | Object Category Update, lihat 12.4 | 200, category yang diperbarui |
| `DELETE /api/categories/{id}` | ID category yang dihapus | Tidak ada | 204, tanpa body response |
| `GET /api/orders` | Tidak ada | Tidak ada | 200, array ringkasan order |
| `GET /api/orders/{id}` | ID order | Tidak ada | 200, order dengan items |
| `POST /api/orders` | Tidak ada | Object Order Create, lihat 12.5 | 201, order dengan items dan total |

Contoh request tanpa body:

```http
GET /health
GET /api/products
GET /api/products/1
GET /api/categories
GET /api/categories/1
GET /api/orders
GET /api/orders/30000
DELETE /api/products/30000
DELETE /api/categories/30000
```

Setiap baris di atas adalah request terpisah. GET detail/DELETE pada ID yang tidak ada menghasilkan 404 dengan kode `PRODUCT_NOT_FOUND`, `CATEGORY_NOT_FOUND`, atau `ORDER_NOT_FOUND` sesuai module. DELETE product yang sudah dipakai order details atau category yang masih dipakai product menghasilkan 409 `FOREIGN_KEY_CONFLICT` pada schema FK standar.

### 12.3. Payload Product

| Field | Tipe JSON | POST | PUT | Validasi/perilaku |
|---|---|---|---|---|
| `product_id` | integer | Wajib | Opsional | POST: 1–32767, unik. PUT: tidak diisi atau 0 memakai ID path; nilai selain 0 harus sama dengan ID path |
| `product_name` | string | Wajib | Wajib | Tidak kosong/whitespace, maksimal 40 karakter |
| `supplier_id` | integer atau `null` | Opsional | Opsional | Jika berupa angka, 1–32767 dan harus ada di suppliers. Tidak diisi menjadi `null` |
| `category_id` | integer atau `null` | Opsional | Opsional | Jika berupa angka, 1–32767 dan harus ada di categories. Tidak diisi menjadi `null` |
| `unit_price` | number | Opsional, default 0 | Opsional, menjadi 0 jika tidak diisi | Tidak negatif dan tidak melebihi `math.MaxFloat32` (sekitar 3.4028235e38) |
| `units_in_stock` | integer | Opsional, default 0 | Opsional, menjadi 0 jika tidak diisi | 0–32767 |
| `discontinued` | integer | Opsional, default 0 | Opsional, menjadi 0 jika tidak diisi | 0 = aktif, 1 = dihentikan |

**POST `/api/products` — buat product**

Buat category 30000 terlebih dahulu bila memakai contoh ini.

```json
{
  "product_id": 30000,
  "product_name": "Produk Workshop",
  "supplier_id": null,
  "category_id": 30000,
  "unit_price": 10,
  "units_in_stock": 20,
  "discontinued": 0
}
```

**PUT `/api/products/30000` — ganti field product**

```json
{
  "product_name": "Produk Workshop Baru",
  "supplier_id": null,
  "category_id": 30000,
  "unit_price": 12.5,
  "units_in_stock": 30,
  "discontinued": 0
}
```

Input tidak valid menghasilkan 400 `INVALID_INPUT`. Create dengan ID yang sudah dipakai menghasilkan 409 `DUPLICATE_KEY`. Supplier/category yang tidak ada menghasilkan 409 `FOREIGN_KEY_CONFLICT`.

### 12.4. Payload Category

| Field | Tipe JSON | POST | PUT | Validasi/perilaku |
|---|---|---|---|---|
| `category_id` | integer | Wajib | Opsional | POST: 1–32767, unik. PUT: tidak diisi atau 0 memakai ID path; nilai selain 0 harus sama dengan ID path |
| `category_name` | string | Wajib | Wajib | Tidak kosong/whitespace, maksimal 15 karakter |
| `description` | string atau `null` | Opsional | Opsional | Tidak diisi menjadi `null`; string kosong disimpan sebagai string kosong |

**POST `/api/categories` — buat category**

```json
{
  "category_id": 30000,
  "category_name": "Workshop",
  "description": "Kategori latihan"
}
```

**PUT `/api/categories/30000` — ganti field category**

```json
{
  "category_name": "Workshop Baru",
  "description": "Kategori latihan diperbarui"
}
```

Input tidak valid menghasilkan 400 `INVALID_INPUT`. Create dengan ID yang sudah dipakai menghasilkan 409 `DUPLICATE_KEY`. Gunakan `"description": null` untuk mengosongkan deskripsi menjadi SQL NULL.

### 12.5. Payload Order

**POST `/api/orders` — buat order dan kurangi stok**

| Field | Tipe JSON | Wajib | Validasi/perilaku |
|---|---|---|---|
| `order_id` | integer | Ya | 1–32767, belum dipakai order lain |
| `customer_id` | string | Ya | Tidak kosong/whitespace, maksimal 5 karakter, harus ada di customers |
| `items` | array object | Ya | Minimal satu item, tidak boleh ada product ID berulang |
| `items[].product_id` | integer | Ya | 1–32767, produk harus ada dan aktif |
| `items[].quantity` | integer | Ya | 1–32767, tidak melebihi stok tersedia |
| `items[].discount` | number | Tidak, default 0 | 0–1; `0.1` berarti diskon 10%, `1` berarti diskon 100% |

```json
{
  "order_id": 30000,
  "customer_id": "ALFKI",
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "discount": 0
    },
    {
      "product_id": 2,
      "quantity": 3,
      "discount": 0.1
    }
  ]
}
```

Customer `ALFKI` dan product 1/2 mengacu dataset Northwind umum; periksa data aktual sebelum menjalankan. `order_date` diisi database saat create. Jangan kirim `unit_price`, `subtotal`, `total`, atau `order_date`: field tersebut bukan bagian payload create dan akan ditolak sebagai field tidak dikenal. Harga berasal dari database, subtotal/total dihitung service.

| Kondisi | Status | Error code |
|---|---|---|
| Field order tidak valid, items kosong, atau product ID berulang | 400 | `INVALID_INPUT` |
| Customer tidak ada | 404 | `CUSTOMER_NOT_FOUND` |
| Product tidak ada | 404 | `PRODUCT_NOT_FOUND` |
| Order ID sudah dipakai | 409 | `DUPLICATE_KEY` |
| Stok kurang atau product discontinued | 409 | `INSUFFICIENT_STOCK` |
| Harga produk NULL/tidak valid | 409 | `INVALID_PRODUCT_PRICE` |
| Konflik foreign key | 409 | `FOREIGN_KEY_CONFLICT` |

Request sukses mengurangi stok dan menyimpan order/details dalam satu transaksi. Error sebelum commit membatalkan seluruh perubahan. Untuk membaca hasil, gunakan `GET /api/orders/30000` tanpa body. API belum menyediakan PUT/DELETE Order.

### 12.6. Contoh curl dan response

Buka `examples/requests.http` memakai VS Code REST Client lalu jalankan setiap request sesuai urutan. Alternatif curl:

```bash
curl http://localhost:8080/api/products/1

curl -i -X POST http://localhost:8080/api/categories \
  -H 'Content-Type: application/json' \
  -d '{"category_id":30000,"category_name":"Workshop","description":"Latihan"}'

curl -i -X POST http://localhost:8080/api/products \
  -H 'Content-Type: application/json' \
  -d '{"product_id":30000,"product_name":"Produk Workshop","category_id":30000,"unit_price":10,"units_in_stock":20,"discontinued":0}'

curl -i -X PUT http://localhost:8080/api/products/30000 \
  -H 'Content-Type: application/json' \
  -d '{"product_name":"Produk Baru","category_id":30000,"unit_price":12.5,"units_in_stock":30,"discontinued":0}'

curl -i -X DELETE http://localhost:8080/api/products/30000
curl -i -X DELETE http://localhost:8080/api/categories/30000

curl -i -X POST http://localhost:8080/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"order_id":30000,"customer_id":"ALFKI","items":[{"product_id":1,"quantity":2,"discount":0},{"product_id":2,"quantity":3,"discount":0.1}]}'

curl http://localhost:8080/api/orders/30000
```

Contoh bentuk response product (nilai mengikuti isi database):

```json
{"data":{"product_id":1,"product_name":"Chai","supplier_id":1,"category_id":1,"unit_price":18,"units_in_stock":39,"discontinued":0}}
```

Contoh error:

```json
{"error":{"code":"PRODUCT_NOT_FOUND","message":"Record not found"}}
```

Validasi gagal/JSON rusak/field tak dikenal: 400. ID atau record hilang: 404 (ID di luar rentang atau bukan angka: 400). Duplicate PK, FK konflik, stok kurang, discontinued, harga produk tidak tersedia: 409. Error tak terduga: 500 tanpa detail SQL di response. DELETE sukses: 204 tanpa body. JSON hanya menerima satu nilai dengan batas body 1 MiB. Harga/stok negatif, quantity tidak positif, discount di luar 0–1, dan item duplikat ditolak.

## 13. Unit test dan pemeriksaan kode

```bash
gofmt -w .
go test ./...
go vet ./...
go build ./...
```

Pemeriksaan tambahan memakai Staticcheck v0.5.1 yang mendukung Go 1.22:

```bash
go install honnef.co/go/tools/cmd/staticcheck@v0.5.1
"$(go env GOPATH)/bin/staticcheck" ./...
```

Perintah tersebut memasang tool lint terpisah, tanpa menambah dependency aplikasi. Jika mengatur `GOBIN`, jalankan binary dari direktori tersebut.

Test tidak memuat `.env` dan tidak menghubungi database. Cakupan: path ID valid/tidak valid, JSON rusak/unknown/trailing, validasi input, method 405, health, envelope JSON/error, pemetaan error PostgreSQL, dan perhitungan diskon/subtotal. Dependency nil digunakan hanya untuk jalur yang harus berhenti sebelum akses database. Tidak ada integration test otomatis; transaksi/locking/FK perlu diverifikasi secara langsung pada database latihan yang schema-nya sesuai.

## 14. Skenario praktikum berurutan

1. Siapkan salinan database latihan. Pastikan category/product/order ID 30000 belum dipakai, customer `ALFKI` ada, product 1 dan 2 aktif dan stok cukup. Jangan menimpa data milik orang lain.
2. Jalankan server dan health check; baca daftar product/category dan product 1.
3. Buat category 30000, baca detail, lalu ubah namanya.
4. Buat product 30000 yang mengacu category tersebut. Baca detail dan ubah harga/stok.
5. Coba harga negatif, JSON dengan field salah, ID tidak valid, dan create dengan ID sama. Amati 400/409.
6. Coba hapus category 30000 saat product masih memakainya. Harapkan 409 pada FK standar.
7. Hapus product latihan lalu category. Harapkan 204; GET ID yang dihapus menghasilkan 404.
8. Catat stok product 1 dan 2. Buat order 30000 lewat contoh request. Baca detail dan cek total; stok masing-masing berkurang 2 dan 3.
9. Kirim ulang order ID sama. Harapkan 409 dan stok tidak berkurang lagi karena rollback. Coba order ID lain dengan quantity melebihi stok atau product berulang; harapkan 409/400 tanpa perubahan parsial.

Pembuatan order mengubah stok nyata. Endpoint hapus order tidak disediakan. Untuk mengulang seluruh praktikum, gunakan salinan database baru atau reset data latihan secara terkontrol lewat SQL setelah memeriksa relasinya.

## 15. Troubleshooting PostgreSQL

| Gejala | Pemeriksaan |
|---|---|
| Startup gagal | Pastikan `.env` berada di working directory, `DATABASE_URL` terisi; periksa error jaringan atau SQLSTATE pada log; error parsing DSN tetap disamarkan agar password tidak bocor |
| DNS/timeout/network unreachable | Periksa host, internet, firewall, project Supabase aktif, dan IPv6; coba Session Pooler untuk IPv4 |
| Password authentication failed | Salin user sesuai mode direct/session, reset password jika perlu, URL-encode password |
| SSL error | Gunakan parameter `sslmode=require` |
| Relation/column does not exist | Jalankan query inspeksi schema; nama tabel/kolom harus cocok pada `public` |
| Boolean/integer mismatch | Varian schema menggunakan `discontinued boolean`; sesuaikan model, SQL, dan pengecekan service |
| Permission denied/RLS | Periksa hak role terhadap tabel schema public dan policy RLS yang berlaku; gunakan role database workshop yang sesuai |
| POST 409 | ID sudah digunakan, FK tidak ada, stok kurang, atau produk discontinued; periksa data sebelum mencoba ulang |
| DELETE category 409 | Hapus/pindahkan product yang masih memakai category terlebih dahulu |
| Terlalu banyak koneksi | Periksa limit project/pooler; kurangi `SetMaxOpenConns` sesuai jatah workshop |

Boilerplate ditujukan untuk latihan lokal: belum menambahkan autentikasi, pagination, atau aturan produksi tambahan. Tambahkan saat kebutuhan aplikasi memerlukannya; jangan membuka API tulis ini ke publik tanpa autentikasi.
