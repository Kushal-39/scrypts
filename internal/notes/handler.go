package notes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"scrypts/internal/auth"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Note struct {
	ID       string
	Owner    string
	Content  []byte
	Nonce    []byte //used for encryption
	Created  int64
	Modified int64
}

var notes = make(map[string]Note)

type NoteReq struct {
	Content string
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req NoteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Extract JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	noteID := uuid.New().String()
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		http.Error(w, "Failed to generate nonce", http.StatusInternalServerError)
		return
	}
	key := []byte("thisis32byteslongpassphraseforaes!")
	block,err:=aes.NewCipher(key)
	if err!=nil{
		http.Error(w, "Failed to create cipher", http.StatusInternalServerError)
		return
	}
	gcm, err:= cipher.NewGCM(block)
	if err!=nil{
		http.Error(w, "Failed to create GCM", http.StatusInternalServerError)
		return
	}
	ciphertext:=gcm.Seal(nil, nonce, []byte(req.Content), nil)

	note := Note{
		ID: noteID,
		Owner: username,
		Content: ciphertext,
		Nonce: nonce,
		Created: time.Now().Unix(),
		Modified: time.Now().Unix(),
	}
	notes[noteID]=note
	json.NewEncoder(w).Encode(map[string]string{"id":noteID})
}
