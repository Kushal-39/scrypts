package main

import (
	"fmt"
	"net/http"
	"scrypts/internal/auth"
	"scrypts/internal/notes"
)

func main() {
	fmt.Println("Starting server on http://localhost:8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Scrypts is alive and kicking")
	})
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodPost:
			notes.CreateNoteHandler(w,r)
		case http.MethodGet:
			notes.GetNotesHandler(w,r)
		case http.MethodPut:
			notes.UpdateNoteHandler(w, r)
		case http.MethodDelete:
			notes.DeleteNoteHandler(w,r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
