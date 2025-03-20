package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// NoteRepository interface defines the methods we need to mock
type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetAllNotes(ctx context.Context) ([]*models.Note, error)
	GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error)
}

// MockNoteRepository is a mock implementation of our repository
type MockNoteRepository struct {
	mock.Mock
}

// CreateNote mocks the CreateNote method
func (m *MockNoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

// GetAllNotes mocks the GetAllNotes method
func (m *MockNoteRepository) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Note), args.Error(1)
}

// GetNoteByID mocks the GetNoteByID method
func (m *MockNoteRepository) GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

// NoteHandler handles Note-related HTTP requests
type NoteHandler struct {
	repo NoteRepository
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(repo NoteRepository) *NoteHandler {
	return &NoteHandler{
		repo: repo,
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
	if err := h.repo.CreateNote(r.Context(), note); err != nil {
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
	notes, err := h.repo.GetAllNotes(r.Context())
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

	note, err := h.repo.GetNoteByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// SetupRouter sets up the router with all the note routes
func SetupRouter(repo NoteRepository) *mux.Router {
	router := mux.NewRouter()
	noteHandler := NewNoteHandler(repo)

	router.HandleFunc("/notes", noteHandler.GetAllNotes).Methods("GET")
	router.HandleFunc("/notes", noteHandler.CreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}", noteHandler.GetNoteByID).Methods("GET")

	return router
}

func TestCreateNote(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Setup expectations
	mockRepo.On("CreateNote", mock.Anything, mock.MatchedBy(func(note *models.Note) bool {
		return note.Title == "Test Note"
	})).Return(nil)

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request
	reqBody := bytes.NewBufferString(`{"title": "Test Note"}`)
	req, _ := http.NewRequest("POST", "/notes", reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse the response
	var note models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &note)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, "Test Note", note.Title)
	assert.NotEmpty(t, note.ID)

	// Verify that the mock method was called
	mockRepo.AssertExpectations(t)
}

func TestCreateNoteWithEmptyTitle(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// No expectations needed, as the handler should return before calling the repo

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request with empty title
	reqBody := bytes.NewBufferString(`{"title": ""}`)
	req, _ := http.NewRequest("POST", "/notes", reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verify that the mock method was not called
	mockRepo.AssertNotCalled(t, "CreateNote")
}

func TestGetAllNotes(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Create test data
	note1 := models.NewNote("Note 1")
	note2 := models.NewNote("Note 2")
	notes := []*models.Note{note1, note2}

	// Setup expectations
	mockRepo.On("GetAllNotes", mock.Anything).Return(notes, nil)

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request
	req, _ := http.NewRequest("GET", "/notes", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var responseNotes []*models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &responseNotes)
	assert.NoError(t, err)

	// Verify the response
	assert.Len(t, responseNotes, 2)
	assert.Equal(t, "Note 1", responseNotes[0].Title)
	assert.Equal(t, "Note 2", responseNotes[1].Title)

	// Verify that the mock method was called
	mockRepo.AssertExpectations(t)
}

func TestGetNoteByID(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Create test data
	note := models.NewNote("Test Note")
	noteID := note.ID

	// Setup expectations
	mockRepo.On("GetNoteByID", mock.Anything, noteID).Return(note, nil)

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request
	url := fmt.Sprintf("/notes/%s", noteID)
	req, _ := http.NewRequest("GET", url, nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var responseNote models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &responseNote)
	assert.NoError(t, err)

	// Verify the response
	assert.Equal(t, noteID, responseNote.ID)
	assert.Equal(t, "Test Note", responseNote.Title)

	// Verify that the mock method was called
	mockRepo.AssertExpectations(t)
}

func TestGetNoteByIDWithInvalidID(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request with invalid ID
	req, _ := http.NewRequest("GET", "/notes/invalid-id", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verify that the mock method was not called
	mockRepo.AssertNotCalled(t, "GetNoteByID")
}

func TestGetNoteByIDNotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Generate a random ID that will not be found
	noteID := uuid.New()

	// Setup expectations
	mockRepo.On("GetNoteByID", mock.Anything, noteID).Return(nil, fmt.Errorf("note not found"))

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request
	url := fmt.Sprintf("/notes/%s", noteID)
	req, _ := http.NewRequest("GET", url, nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify that the mock method was called
	mockRepo.AssertExpectations(t)
}

func TestCreateNoteWithInvalidJSON(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// No expectations needed, as the handler should return before calling the repo

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request with invalid JSON
	reqBody := bytes.NewBufferString(`{"title": 123}`) // Title should be a string, not a number
	req, _ := http.NewRequest("POST", "/notes", reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verify that the mock method was not called
	mockRepo.AssertNotCalled(t, "CreateNote")
}

func TestGetAllNotesWithError(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockNoteRepository)

	// Setup expectations to return an error
	mockRepo.On("GetAllNotes", mock.Anything).Return(nil, fmt.Errorf("database error"))

	// Setup router with mock repository
	router := SetupRouter(mockRepo)

	// Create a test request
	req, _ := http.NewRequest("GET", "/notes", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Handle the request
	router.ServeHTTP(rr, req)

	// Check the status code - should be server error
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify that the mock method was called
	mockRepo.AssertExpectations(t)
}
