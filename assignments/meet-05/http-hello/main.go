package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type userRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func myFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		x, y := 10, 12
		fmt.Fprintln(w, "Nilai x:", x, "\ndan nilai y:", y, "\nhasilnya x*y:", x*y)
	case http.MethodPost:
		var userReq userRequest
		err := json.NewDecoder(r.Body).Decode(&userReq)
		if err != nil {
			http.Error(w, "Json Gagal di decode!", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "User Request: %+v\n", userReq)
	case http.MethodPut:
		userId := r.URL.Query().Get("id")
		fmt.Fprintln(w, "User ID yang akan diupdate adalah :", userId)
	case http.MethodDelete:
		userId := r.URL.Query().Get("id")
		fmt.Fprintln(w, "User ID yang akan dihapus adalah :", userId)
	default:
		http.Error(w, "Method tidak valid", http.StatusMethodNotAllowed)
	}

}

func main() {
	//root request handler -> user merequest data ke server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go HTTP service")
	})

	http.HandleFunc("/myfunc", myFunc)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok kamu di /healthz")
	})

	log.Println("server listening on :8024")
	log.Fatal(http.ListenAndServe(":8024", nil))
}
