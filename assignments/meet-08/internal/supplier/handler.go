package supplier

import (
	"net/http"
)

// Handler exposes the supplier HTTP endpoint.
type Handler struct{ Service Service }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/suppliers" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method tidak didukung", http.StatusMethodNotAllowed)
		return
	}
	if h.Service.Repository.DB == nil {
		http.Error(w, "database belum siap", http.StatusInternalServerError)
		return
	}
	rows, err := h.Service.List(r.Context())
	if err != nil {
		http.Error(w, "query database gagal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rows)
}
