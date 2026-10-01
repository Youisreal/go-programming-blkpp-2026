package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//root request handler -> user merequest data ke server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go HTTP service")
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok kamu di /healthz")
	})

	log.Println("server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
