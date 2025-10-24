package notes

import (
	"encoding/json"
	"log"
	"net/http"
	"scrypts/internal/auth"
	"scrypts/internal/storage"
	"scrypts/internal/utils"
	"time"

	"github.com/google/uuid"
)

type NoteReq struct {
	Content string `json:"content"`
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

	if len(req.Content) > storage.MaxNoteContentSize {
		http.Error(w, "Note content too large", http.StatusRequestEntityTooLarge)
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
		log.Printf("EncryptAESGCM error : %v", err)
		http.Error(w, "Failed to encrypt note", http.StatusInternalServerError)
		return
	}

	now := time.Now().Unix()
	snote := storage.Note{
		ID:       noteID,
		Owner:    username,
		Content:  ciphertext,
		Nonce:    nonce,
		Created:  now,
		Modified: now,
	}
	if err := storage.SaveNote(snote); err != nil {
		log.Printf("Savenote error: %v", err)
		http.Error(w, "Failed to Save Note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
	snotes, err := storage.GetNotesByOwner(username)
	if err != nil {
		log.Printf("GetNotesByOwner error: %v", err)
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}

	userKey, err := auth.GetUserKey(username)
	if err != nil {
		log.Printf("Get User Key error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	type NoteResp struct {
		ID       string `json:"id"`
		Owner    string `json:"owner"`
		Content  string `json:"content"`
		Created  int64  `json:"created"`
		Modified int64  `json:"modified"`
	}
	resp :=make([]NoteResp,0,len(snotes))
	for _, sn := range snotes{
		pt, derr:=utils.DecryptAESGCM(userKey,sn.Nonce,sn.Content)
		if derr!=nil{
			log.Printf("DecryptAESGCM error :%v",err)
			http.Error(w,"failed to decrypt note", http.StatusInternalServerError)
			return
		}
		resp = append(resp, NoteResp{
			ID: sn.ID,
			Owner: sn.Owner,
			Content: string(pt),
			Created: sn.Created,
			Modified: sn.Modified,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	// validate ID and content early
	if _, err := uuid.Parse(req.ID); err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}
	if len(req.Content) > storage.MaxNoteContentSize {
		http.Error(w, "Note content too large", http.StatusRequestEntityTooLarge)
		return
	}
	snotes, err := storage.GetNotesByOwner(username)
	if err != nil {
		log.Printf("GetNotesByOwner err: %v", err)
		http.Error(w, "failed to query notes", http.StatusInternalServerError)
		return
	}
	var existing *storage.Note
	for i := range snotes {
		if snotes[i].ID == req.ID {
			existing = &snotes[i]
			break
		}
	}
	if existing == nil {
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
	now := time.Now().Unix()
	snote := storage.Note{
		ID:       req.ID,
		Owner:    username,
		Content:  ciphertext,
		Nonce:    nonce,
		Created:  existing.Created,
		Modified: now,
	}
	if err := storage.UpdateNote(snote); err != nil {
		log.Printf("UpdateNote error: %v", err)
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}
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
	if _, err := uuid.Parse(req.ID); err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}
	snotes, err := storage.GetNotesByOwner(username)
	if err != nil {
		log.Printf("GetNotesByOwner error: %v", err)
		http.Error(w, "failed to query note", http.StatusInternalServerError)
		return
	}
	found := false
	for i := range snotes {
		if snotes[i].ID == req.ID {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err := storage.DeleteNote(req.ID, username); err != nil {
		log.Printf("DeleteNote error: %v", err)
		http.Error(w, "failed to delete note", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
