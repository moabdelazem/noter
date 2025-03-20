package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/database"
	"github.com/moabdelazem/noter/internal/models"
)

// NoteHandler handles HTTP requests for notes
type NoteHandler struct {
	noteRepo *database.NoteRepository
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(noteRepo *database.NoteRepository) *NoteHandler {
	return &NoteHandler{
		noteRepo: noteRepo,
	}
}

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Title string `json:"title"`
}

// CreateNote handles the request to create a new note
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	note := models.NewNote(req.Title)
	if err := h.noteRepo.CreateNote(r.Context(), note); err != nil {
		errMsg := fmt.Sprintf("Failed to create note: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetAllNotes handles the request to get all notes
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteRepo.GetAllNotes(r.Context())
	if err != nil {
		http.Error(w, "Failed to get notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// GetNoteByID handles the request to get a note by ID
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.noteRepo.GetNoteByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}
