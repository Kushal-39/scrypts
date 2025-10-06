package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make(map[string]string)

func isComplex(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _,c := range password{
		switch{
		case unicode.IsUpper(c):
			hasUpper=true
		case unicode.IsLower(c):
			hasLower=true
		case unicode.IsDigit(c):
			hasDigit=true
		case unicode.IsPunct(c):
			hasSpecial=true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Username) < 4 || len(req.Password) < 8 {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	if !isComplex(req.Password){
		http.Error(w, "Password is weak(must contain uppercase,lowercase,digit,symbol)",http.StatusBadRequest)
		return
	}
	if _, exists := users[req.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hashed, err := HashPass(req.Password)
	if err != nil {
		http.Error(w, "Error in hashing password", http.StatusInternalServerError)
		return
	}

	users[req.Username] = hashed
	fmt.Fprintln(w, "User registered successfully")
}
