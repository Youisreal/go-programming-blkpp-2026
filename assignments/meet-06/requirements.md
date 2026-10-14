# Requirements Meet 06

## Wajib

- Lengkapi semua TODO di `starter/meet-06`.
- Gunakan Go standard library saja.
- Gunakan `int64` untuk semua nilai uang.
- Validasi input sebelum perhitungan.
- Return error, jangan panic, untuk data transaksi yang salah.
- Pertahankan CLI berbasis file JSON.

## Function yang Harus Selesai

```go
func CalculateLineTotal(price int64, qty int64) (int64, error)
func CalculateDiscount(subtotal int64, percent int64) (int64, error)
func CalculateTax(amount int64, percent int64) (int64, error)
func CalculateChange(grandTotal int64, amountPaid int64) (int64, error)
func CalculateInvoice(request InvoiceRequest) (InvoiceResult, error)
```

Function validasi boleh dipecah jika membantu, tetapi jangan tambah abstraksi besar.

## Larangan

- Jangan pakai `float32` atau `float64` untuk uang.
- Jangan tambah HTTP endpoint.
- Jangan tambah database.
- Jangan tambah framework.
- Jangan hardcode hasil khusus untuk data contoh.

## Output

Output harus berupa JSON dengan field:

- `invoice_number`
- `customer_name`
- `items`
- `summary`

`summary.status` bernilai `"paid"` jika pembayaran cukup.
