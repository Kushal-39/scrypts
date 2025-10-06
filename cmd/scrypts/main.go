package main

import (
	"fmt"
	"net/http"
	"scrypts/internal/auth"
)

func main() {
	fmt.Println("Starting server on http://localhost:8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Scrypts is alive and kicking")
	})
	http.HandleFunc("/register", auth.RegisterHandler)
	http.ListenAndServe(":8080", nil)
}
