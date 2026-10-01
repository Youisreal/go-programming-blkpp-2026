package main

import (
	"fmt"
	"log"
	"net/http"
)

func myFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		x, y := 10, 12
		fmt.Fprintln(w, "Nilai x:", x, "\ndan nilai y:", y, "\nhasilnya x*y:", x*y)
	} else if r.Method == http.MethodPost {
		fmt.Fprintln(w, "Ini adalah Method POST")
	} else if r.Method == http.MethodPut {
		fmt.Fprintln(w, "Ini adalah Method Update")
	} else if r.Method == http.MethodDelete {
		fmt.Fprintln(w, "Ini adalah Method Delete")
	} else {
		fmt.Fprintln(w, "Method tidak valid")
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
