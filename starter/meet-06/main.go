package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type errorResponse struct {
	Error string `json:"error"`
}

type store struct {
	mu       sync.RWMutex
	file     string
	products map[string]Product
	stock    map[string]int64
	invoices map[string]InvoiceResult
}

func newStore(file string) (*store, error) {
	store := &store{
		file:     file,
		products: map[string]Product{},
		stock:    map[string]int64{},
		invoices: map[string]InvoiceResult{},
	}
	return store, store.loadProducts()
}

func main() {
	store, err := newStore("data/products.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", newMux(store)))
}

func newMux(store *store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", store.listProducts)
	mux.HandleFunc("POST /products", store.createProduct)
	mux.HandleFunc("GET /products/{sku}", store.getProduct)
	mux.HandleFunc("GET /stock", store.listStock)
	mux.HandleFunc("POST /stock", store.adjustStock)
	mux.HandleFunc("POST /checkout", store.checkout)
	mux.HandleFunc("GET /invoices", store.listInvoices)
	mux.HandleFunc("GET /invoices/{invoice_number}", store.getInvoice)
	return mux
}

func (s *store) listProducts(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *store) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if !decodeJSON(w, r, &product) {
		return
	}
	if err := ValidateProduct(product); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.products[product.SKU]; exists {
		writeError(w, http.StatusConflict, "sku sudah ada")
		return
	}
	s.products[product.SKU] = product
	s.stock[product.SKU] = 0
	if err := s.saveProducts(); err != nil {
		delete(s.products, product.SKU)
		delete(s.stock, product.SKU)
		writeError(w, http.StatusInternalServerError, "gagal menyimpan product")
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

func (s *store) getProduct(w http.ResponseWriter, r *http.Request) {
	sku := r.PathValue("sku")

	s.mu.RLock()
	product, exists := s.products[sku]
	s.mu.RUnlock()
	if !exists {
		writeError(w, http.StatusNotFound, "product tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *store) listStock(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]StockItem, 0, len(s.products))
	for sku, product := range s.products {
		items = append(items, StockItem{
			SKU:      sku,
			Name:     product.Name,
			Quantity: s.stock[sku],
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *store) adjustStock(w http.ResponseWriter, r *http.Request) {
	var adjustment StockAdjustment
	if !decodeJSON(w, r, &adjustment) {
		return
	}
	if err := ValidateStockAdjustment(adjustment); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	product, exists := s.products[adjustment.SKU]
	if !exists {
		writeError(w, http.StatusNotFound, "product tidak ditemukan")
		return
	}
	current := s.stock[adjustment.SKU]
	if adjustment.Quantity > 0 && current > math.MaxInt64-adjustment.Quantity {
		writeError(w, http.StatusBadRequest, "stock overflow")
		return
	}
	next := current + adjustment.Quantity
	if next < 0 {
		writeError(w, http.StatusBadRequest, "stock tidak boleh negatif")
		return
	}
	s.stock[adjustment.SKU] = next
	writeJSON(w, http.StatusOK, StockItem{SKU: adjustment.SKU, Name: product.Name, Quantity: next})
}

func (s *store) checkout(w http.ResponseWriter, r *http.Request) {
	var checkout CheckoutRequest
	if !decodeJSON(w, r, &checkout) {
		return
	}
	if err := ValidateCheckoutRequest(checkout); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.invoices[checkout.InvoiceNumber]; exists {
		writeError(w, http.StatusConflict, "invoice_number sudah ada")
		return
	}

	request, err := s.buildInvoiceRequest(checkout)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := CalculateInvoice(request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	for _, item := range checkout.Items {
		s.stock[item.SKU] -= item.Quantity
	}
	s.invoices[result.InvoiceNumber] = result
	writeJSON(w, http.StatusCreated, result)
}

func (s *store) buildInvoiceRequest(checkout CheckoutRequest) (InvoiceRequest, error) {
	items := make([]TransactionItem, 0, len(checkout.Items))
	requested := map[string]int64{}
	for _, item := range checkout.Items {
		product, exists := s.products[item.SKU]
		if !exists {
			return InvoiceRequest{}, fmt.Errorf("product %s tidak ditemukan", item.SKU)
		}
		if !product.Active {
			return InvoiceRequest{}, fmt.Errorf("product %s tidak aktif", item.SKU)
		}
		if requested[item.SKU] > math.MaxInt64-item.Quantity {
			return InvoiceRequest{}, fmt.Errorf("quantity %s overflow", item.SKU)
		}
		requested[item.SKU] += item.Quantity
		if s.stock[item.SKU] < requested[item.SKU] {
			return InvoiceRequest{}, fmt.Errorf("stock %s tidak cukup", item.SKU)
		}
		items = append(items, TransactionItem{
			SKU:       product.SKU,
			Name:      product.Name,
			UnitPrice: product.Price,
			Quantity:  item.Quantity,
		})
	}
	return InvoiceRequest{
		InvoiceNumber:   checkout.InvoiceNumber,
		Customer:        checkout.Customer,
		Items:           items,
		DiscountPercent: checkout.DiscountPercent,
		TaxPercent:      checkout.TaxPercent,
		ShippingFee:     checkout.ShippingFee,
		AmountPaid:      checkout.AmountPaid,
	}, nil
}

func (s *store) listInvoices(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]InvoiceResult, 0, len(s.invoices))
	for _, invoice := range s.invoices {
		results = append(results, invoice)
	}
	writeJSON(w, http.StatusOK, results)
}

func (s *store) getInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := r.PathValue("invoice_number")

	s.mu.RLock()
	result, exists := s.invoices[invoiceNumber]
	s.mu.RUnlock()
	if !exists {
		writeError(w, http.StatusNotFound, "invoice tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, data any) bool {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(data); err != nil {
		writeError(w, http.StatusBadRequest, "JSON request tidak valid: "+err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("write response:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func (s *store) loadProducts() error {
	file, err := os.Open(s.file)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var products []Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return err
	}
	for _, product := range products {
		if err := ValidateProduct(product); err != nil {
			return err
		}
		s.products[product.SKU] = product
		s.stock[product.SKU] = 0
	}
	return nil
}

func (s *store) saveProducts() error {
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return err
	}

	products := make([]Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}

	temp := s.file + ".tmp"
	file, err := os.Create(temp)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temp, s.file)
}
