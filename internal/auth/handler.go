package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"unicode"
	"time"
    "github.com/golang-jwt/jwt/v5"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var users = make(map[string]string)

func generateJWT(username string) (string,error){
	claims:= jwt.MapClaims{
		"username":username,
		"exp": time.Now().Add(15*time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(jwtSecret)
}

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

func LoginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err!=nil{
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	hashed,exists := users[req.Username]
	if !exists{
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if !CheckPasswordHash(req.Password, hashed){
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(req.Username)
	if err != nil{
		http.Error(w, "Could not generate Token", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}