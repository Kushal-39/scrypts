package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"net/http"
	"scrypts/internal/config"
	"scrypts/internal/storage"
	"scrypts/internal/utils"
	"strings"
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
		case !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c):
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
		http.Error(w, "Password must contain uppercase, lowercase, digit, and special character", http.StatusBadRequest)
		return
	}

	// Check if user exists (prevent user enumeration with timing attack mitigation)
	_, err := storage.GetUser(req.Username)
	if err == nil {
		// User exists - add random delay to mimic registration time
		time.Sleep(time.Duration(50+mrand.Intn(50)) * time.Millisecond)
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	hashed, err := HashPass(req.Password)
	if err != nil {
		http.Error(w, "Error in hashing password", http.StatusInternalServerError)
		return
	}

	u := storage.User{
		Username:     req.Username,
		PasswordHash: hashed,
		WrappedKey:   nil,
		WrappedNonce: nil,
		CreatedAt:    time.Now().Unix(),
	}
	if err := storage.CreateUser(u); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	userKey := make([]byte, 32)
	if _, err := rand.Read(userKey); err != nil {
		fmt.Println("warning: failed to generate user encryption key:", err)
	} else {
		wrapped, nonce, werr := utils.WrapKey(config.MasterKey, userKey)
		if werr != nil {
			fmt.Println("warning: failed to wrap user key:", werr)
		} else {
			if err := storage.SaveWrappedKey(req.Username, wrapped, nonce); err != nil {
				fmt.Println("warning: failed to save wrapped key to DB:", err)
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
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

	// Always perform timing-consistent operations
	u, err := storage.GetUser(req.Username)

	// Dummy hash for timing attack mitigation when user doesn't exist
	dummyHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	hashToCheck := dummyHash
	userValid := false

	if err == nil {
		hashToCheck = u.PasswordHash
		userValid = true
	}

	// Always check password hash (constant time)
	passwordValid := CheckPasswordHash(req.Password, hashToCheck)

	// Only succeed if both user exists and password is correct
	if !userValid || !passwordValid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
	u, err := storage.GetUser(username)
	if err != nil {
		return nil, err
	}
	if len(u.WrappedKey) == 0 || len(u.WrappedNonce) == 0 {
		return nil, fmt.Errorf("no encryption key for user")
	}
	k, err := utils.UnwrapKey(config.MasterKey, u.WrappedNonce, u.WrappedKey)
	if err != nil {
		return nil, err
	}
	return k, nil
}
