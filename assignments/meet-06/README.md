# Assignment Meet 06 - Sales & Invoice Calculator

Tugas ini melanjutkan starter `starter/meet-06`. Lengkapi TODO sampai aplikasi bisa menghitung invoice dari file JSON dan semua test lulus.

## Target

Buat program CLI kasir yang:

- Membaca JSON invoice.
- Memvalidasi data transaksi.
- Menghitung item, subtotal, diskon, pajak simulasi, ongkir, grand total, dan kembalian.
- Mencetak output JSON terstruktur.
- Memakai standard library Go saja.

## Level 1 - Basic

Implementasikan:

- Total per item.
- Subtotal transaksi.
- Final total tanpa diskon dan pajak.

## Level 2 - Intermediate

Implementasikan:

- Diskon.
- Pajak simulasi.
- Ongkir.
- Pembayaran pelanggan.
- Kembalian.
- Validasi input.

## Level 3 - Advanced / Optional

Implementasikan salah satu atau lebih:

- Diskon member berdasarkan `customer.type`.
- Batas maksimum diskon.
- Voucher nominal tetap.
- Alternatif aturan diskon.
- Receipt teks yang mudah dibaca manusia.

Fitur opsional tidak menjadi prasyarat Meet 07.

## Urutan Perhitungan Wajib

1. Hitung `line_total`.
2. Jumlahkan `subtotal`.
3. Hitung `discount_amount`.
4. Hitung `after_discount`.
5. Hitung `tax_amount`.
6. Hitung `grand_total`.
7. Hitung `change`.
8. Isi `status`.

Urutan ini wajib supaya hasil semua peserta deterministik.

## Cara Mengerjakan

```bash
cd starter/meet-06
go test ./...
go run . data/invoice-valid.json
go fmt ./...
```

Test awal akan gagal karena function masih TODO.

## Workflow GitHub

1. Sync fork dengan upstream repository.
2. Buat branch:

```bash
git checkout -b feat/meet-06-invoice-calculator
```

3. Lengkapi starter code di `starter/meet-06`.
4. Jalankan format dan test.
5. Commit dengan pesan deskriptif, contoh:

```bash
git add starter/meet-06 assignments/meet-06
git commit -m "Complete meet 06 invoice calculator"
```

6. Push branch:

```bash
git push origin feat/meet-06-invoice-calculator
```

7. Buka pull request ke branch `main` di fork sendiri, kecuali instruktur menentukan repository submission lain.
8. Tambahkan screenshot atau output terminal yang menunjukkan perhitungan berhasil.

## Acceptance Criteria

- `go test ./...` lulus dari folder `starter/meet-06`.
- `go run . data/invoice-valid.json` mencetak invoice sesuai expected output.
- Semua validasi input utama tersedia.
- Tidak ada panic untuk input user yang salah.
- Tidak memakai framework, database, HTTP server, atau dependency pihak ketiga.
- Uang memakai `int64`, bukan floating point.
- Business logic tidak ditaruh langsung di `main()`.

## Checklist Code Review

- [ ] Function kecil dan jelas.
- [ ] Error message menjelaskan masalah input.
- [ ] Perhitungan persen memakai pembulatan terdekat.
- [ ] Overflow dicegah sebelum perkalian atau penjumlahan.
- [ ] Test lulus.
- [ ] Output JSON rapi dan sesuai struct tag.
- [ ] Challenge opsional tidak merusak behavior utama.
