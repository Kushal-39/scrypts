package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scrypts/internal/auth"
	"scrypts/internal/utils"
	"time"

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

	username, err := auth.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	noteID := uuid.New().String()
	userKey, err := auth.GetUserKey(username)
	if err != nil {
		http.Error(w, "Server error: missing encryption key", http.StatusInternalServerError)
		return
	}
	key := userKey
	nonce, err := utils.GenerateNonce(12)
	if err != nil {
		http.Error(w, "Failed to generate nonce", http.StatusInternalServerError)
		return
	}
	ciphertext, err := utils.EncryptAESGCM(key, nonce, []byte(req.Content))
	if err != nil {
		fmt.Println("EncryptAESGCM error:", err)
		http.Error(w, "Failed to encrypt note", http.StatusInternalServerError)
		return
	}

	note := Note{
		ID:       noteID,
		Owner:    username,
		Content:  ciphertext,
		Nonce:    nonce,
		Created:  time.Now().Unix(),
		Modified: time.Now().Unix(),
	}
	notes[noteID] = note
	json.NewEncoder(w).Encode(map[string]string{"id": noteID})
}

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username, err := auth.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
		return
	}
	userNotes := []Note{}
	for _, note := range notes {
		if note.Owner == username {
			userNotes = append(userNotes, note)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userNotes)
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username, err := auth.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorised access", http.StatusUnauthorized)
		return
	}
	var req struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	note, exists := notes[req.ID]
	if !exists || note.Owner != username {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	userKey, err := auth.GetUserKey(username)
	if err != nil {
		http.Error(w, "Server error: missing encryption key", http.StatusInternalServerError)
		return
	}
	key := userKey
	nonce, err := utils.GenerateNonce(12)
	if err != nil {
		http.Error(w, "Failed to generate Nonce", http.StatusInternalServerError)
		return
	}
	ciphertext, err := utils.EncryptAESGCM(key, nonce, []byte(req.Content))
	if err != nil {
		http.Error(w, "Failed to encrypt note", http.StatusInternalServerError)
		return
	}
	note.Content = ciphertext
	note.Nonce = nonce
	note.Modified = time.Now().Unix()
	notes[req.ID] = note
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username, err := auth.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	note, exists := notes[req.ID]
	if !exists || note.Owner != username {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	delete(notes, req.ID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
