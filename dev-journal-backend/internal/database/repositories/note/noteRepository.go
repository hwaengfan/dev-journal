package noteRepository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	noteModel "github.com/hwaengfan/dev-journal-backend/internal/models/note"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{database: database}
}

// CreateNote creates a new note
func (store *Store) CreateNote(note noteModel.Note) (uuid.UUID, error) {
	noteID := uuid.New()
	query := "INSERT INTO notes (id, userID, linkedProjectID, title, content, favorited, tags) VALUES (?, ?, ?, ?, ?, ?, ?)"

	// convert []string to JSON
	tagsJSON, err := json.Marshal(note.Tags)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert tags to JSON: %v", err)
	}

	_, error := store.database.Exec(query, noteID, note.UserID, note.LinkedProjectID, note.Title, note.Content, note.Favorited, tagsJSON)
	if error != nil {
		return uuid.Nil, fmt.Errorf("failed to create note: %v", error)
	}

	return noteID, nil
}

// GetNotesByLinkedProjectID retrieves all notes by a linked project's ID
func (store *Store) GetNotesByLinkedProjectID(linkedProjectID uuid.UUID) ([]*noteModel.Note, error) {
	// query notes by linked project ID
	query := "SELECT id, linkedProjectID, title, content, favorited, tags, dateCreated, lastEdited FROM notes WHERE linkedProjectID = ?"
	rows, error := store.database.Query(query, linkedProjectID)
	if error != nil {
		return nil, fmt.Errorf("failed to get notes by linked project ID: %v", error)
	}
	defer rows.Close()

	// scan notes from rows
	notes, error := scanNotesFromRows(rows)
	if error != nil {
		return nil, error
	}

	return notes, nil
}

// GetNoteByID retrieves a note by its ID
func (store *Store) GetNoteByID(id uuid.UUID) (*noteModel.Note, error) {
	// query note by ID
	query := "SELECT id, linkedProjectID, title, content, favorited, tags, dateCreated, lastEdited FROM notes WHERE id = ?"
	row := store.database.QueryRow(query, id)

	// scan note from row
	note, error := scanNoteFromRow(row)
	if error != nil {
		return nil, fmt.Errorf("failed to scan note from row: %v", error)
	}

	return note, nil
}

// UpdateNoteByID updates a note by its ID
func (store *Store) UpdateNoteByID(note noteModel.Note, id uuid.UUID) error {
	// base query
	query := "UPDATE notes SET"
	var updates []string
	var args []interface{}

	// conditionally add fields to update
	if note.LinkedProjectID != uuid.Nil {
		updates = append(updates, "linkedProjectID = ?")
		args = append(args, note.LinkedProjectID)
	}
	if note.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, note.Title)
	}
	if note.Content != "" {
		updates = append(updates, "content = ?")
		args = append(args, note.Content)
	}
	if note.Favorited != "" {
		updates = append(updates, "favorited = ?")
		args = append(args, note.Favorited)
	}
	if note.Tags != nil {
		// convert []string to JSON
		tagsJSON, err := json.Marshal(note.Tags)
		if err != nil {
			return fmt.Errorf("failed to convert tags to JSON: %v", err)
		}

		updates = append(updates, "tags = ?")
		args = append(args, tagsJSON)
	}

	// check if there are fields to update
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// finalize query
	query += " " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	// execute the query
	_, error := store.database.Exec(query, args...)
	if error != nil {
		return fmt.Errorf("failed to update note: %v", error)
	}

	return nil
}

// DeleteNoteByID deletes a note by its ID
func (store *Store) DeleteNoteByID(id uuid.UUID) error {
	query := "DELETE FROM notes WHERE id = ?"
	_, error := store.database.Exec(query, id)
	if error != nil {
		return fmt.Errorf("failed to delete note: %v", error)
	}

	return nil
}

// DeleteNotesByLinkedProjectID deletes all notes by a linked project's ID
func (store *Store) DeleteNotesByLinkedProjectID(linkedProjectID uuid.UUID) error {
	query := "DELETE FROM notes WHERE linkedProjectID = ?"
	_, error := store.database.Exec(query, linkedProjectID)
	if error != nil {
		return fmt.Errorf("failed to delete all notes by linked project ID: %v", error)
	}

	return nil
}

func scanNotesFromRows(rows *sql.Rows) ([]*noteModel.Note, error) {
	notes := make([]*noteModel.Note, 0)

	for rows.Next() {
		note := new(noteModel.Note)
		var tagsJSONString string

		error := rows.Scan(&note.ID, &note.LinkedProjectID, &note.Title, &note.Content, &note.Favorited, &tagsJSONString, &note.DateCreated, &note.LastEdited)
		if error != nil {
			return nil, fmt.Errorf("failed to scan note from rows: %v", error)
		}

		// convert JSON to []string
		error = json.Unmarshal([]byte(tagsJSONString), &note.Tags)
		if error != nil {
			return nil, fmt.Errorf("failed to convert tags from JSON: %v", error)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func scanNoteFromRow(row *sql.Row) (*noteModel.Note, error) {
	note := new(noteModel.Note)
	var tagsJSONString string
	error := row.Scan(&note.ID, &note.LinkedProjectID, &note.Title, &note.Content, &note.Favorited, &tagsJSONString, &note.DateCreated, &note.LastEdited)

	if error == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found")
	} else if error != nil {
		return nil, fmt.Errorf("failed to scan note from row: %v", error)
	}

	// convert JSON to []string
	error = json.Unmarshal([]byte(tagsJSONString), &note.Tags)
	if error != nil {
		return nil, fmt.Errorf("failed to convert tags from JSON: %v", error)
	}

	return note, nil
}
