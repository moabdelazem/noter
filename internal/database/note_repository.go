package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/moabdelazem/noter/internal/models"
)

// NoteRepository handles database operations for notes
type NoteRepository struct {
	db *DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *DB) *NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

// CreateNote inserts a new note into the database
func (r *NoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Pool.Exec(ctx, query, note.ID, note.Title, note.CreatedAt, note.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	return nil
}

// GetAllNotes retrieves all notes from the database
func (r *NoteRepository) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notes, nil
}

// GetNoteByID retrieves a note by its ID
func (r *NoteRepository) GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM notes
		WHERE id = $1
	`
	var note models.Note
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&note.ID, &note.Title, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get note by ID: %w", err)
	}
	return &note, nil
}
