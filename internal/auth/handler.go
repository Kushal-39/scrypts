package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"scrypts/internal/config"
	"scrypts/internal/utils"
	"strings"
	"sync"
	"time"
	"unicode"

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

var JwtSecret = config.JwtSecret
var users = make(map[string]string)
var wrappedKeys = make(map[string]struct {
	Wrapped []byte
	Nonce   []byte
})
var wrappedKeysMu sync.RWMutex

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func isComplex(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c):
			hasSpecial = true
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
	if !isComplex(req.Password) {
		http.Error(w, "Password is weak(must contain uppercase,lowercase,digit,symbol)", http.StatusBadRequest)
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
	userKey := make([]byte, 32)
	if _, err := rand.Read(userKey); err != nil {
		fmt.Println("warning: failed to generate user encryption key:", err)
	} else {
		// use configured master key (already 32 bytes)
		wrapped, nonce, werr := utils.WrapKey(config.MasterKey, userKey)
		if werr != nil {
			fmt.Println("warning: failed to wrap user key:", werr)
		} else {
			wrappedKeysMu.Lock()
			wrappedKeys[req.Username] = struct {
				Wrapped []byte
				Nonce   []byte
			}{Wrapped: wrapped, Nonce: nonce}
			wrappedKeysMu.Unlock()
		}
	}
	fmt.Fprintln(w, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	hashed, exists := users[req.Username]
	if !exists {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if !CheckPasswordHash(req.Password, hashed) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(req.Username)
	if err != nil {
		http.Error(w, "Could not generate Token", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetUsernameFromJWT extracts the username claim from a Bearer JWT in the request.
func GetUsernameFromJWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http.ErrNoCookie
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// enforce HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token")
	}
	return username, nil
}

func GetUserKey(username string) ([]byte, error) {
	wrappedKeysMu.RLock()
	entry, ok := wrappedKeys[username]
	wrappedKeysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no encryption key for user")
	}
	k, err := utils.UnwrapKey(config.MasterKey, entry.Nonce, entry.Wrapped)
	if err != nil {
		return nil, err
	}
	return k, nil
}
