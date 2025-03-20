package models

import (
	"time"

	"github.com/google/uuid"
)

// Note represents a note in the database
type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewNote creates a new note with the given title
func NewNote(title string) *Note {
	now := time.Now()
	return &Note{
		ID:        uuid.New(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
