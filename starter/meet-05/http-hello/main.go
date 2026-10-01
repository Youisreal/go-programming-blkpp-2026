package main

import (
	"fmt"
	"log"
	"io"
	"encoding/json"
	"net/http"
)


type UserRequest struct {
	ID		int 		`json:"id"`
	Name 	string 	`json:"name"`
	Email string 	`json:"email"`
}


func HandleJSON(w http.ResponseWriter, r *http.Request) {
  // 1. Baca data mentah menjadi byte
  bodyBytes, err := io.ReadAll(r.Body)
  if err != nil {
		http.Error(w, "Gagal membaca body", http.StatusBadRequest)
    return
  }

  // 2. Cetak string mentah JSON ke terminal
  fmt.Println("--- Mentah JSON String ---")
	fmt.Println(string(bodyBytes))

  w.Write([]byte("JSON berhasil dicetak"))
}



func Myfunc(w http.ResponseWriter, r *http.Request) {
	if(r.Method == http.MethodGet) {
		x := 10
		y := 12 
		fmt.Fprintln(w, 
		"nilai x:", x, 
		"\ndan y:", y, 
		"\nx * y:", x * y)
	} else if(r.Method == http.MethodPost) {
		var userReq UserRequest

		err := json.NewDecoder(r.Body).Decode(&userReq)
		if err != nil {
			fmt.Fprintln(w, "Json Gagal didecode!")
		}

    fmt.Fprintf(w, "Data JSON diterima: %+v\n", userReq)


	} else if(r.Method == http.MethodPut) {
		
		user_id := r.URL.Query().Get("user_id")
		fmt.Fprintln(w, "User ID yang akan diupdate adalah:", user_id)

	} else if(r.Method == http.MethodDelete) {
		
		user_id := r.URL.Query().Get("user_id")
		fmt.Fprintln(w, "User ID yang akan dihapus adalah:", user_id)

	} else {
		fmt.Fprintln(w, "Method tidak valid!")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go HTTP service")
	})

	http.HandleFunc("/myfunc", Myfunc)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	log.Println("server listening on :8000")
	log.Fatal(http.ListenAndServe(":1234", nil))
}
