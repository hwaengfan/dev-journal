package noteModel

import (
	"github.com/google/uuid"
)

type Note struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"userID"`
	LinkedProjectID uuid.UUID `json:"linkedProjectID"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Favorited       string    `json:"favorited"`
	Tags            []string  `json:"tags"`
	DateCreated     string    `json:"dateCreated"`
	LastEdited      string    `json:"lastEdited"`
}

type NoteStore interface {
	CreateNote(note Note) (uuid.UUID, error)
	GetNotesByLinkedProjectID(linkedProjectID uuid.UUID) ([]*Note, error)
	GetNoteByID(id uuid.UUID) (*Note, error)
	UpdateNoteByID(note Note, id uuid.UUID) error
	DeleteNoteByID(id uuid.UUID) error
	DeleteNotesByLinkedProjectID(linkedProjectID uuid.UUID) error
}

type CreateNotePayload struct {
	LinkedProjectID uuid.UUID `json:"linkedProjectID" validate:"required"`
	Title           string    `json:"title" validate:"required"`
	Content         string    `json:"content" validate:"required"`
	Favorited       string    `json:"favorited" validate:"required"`
	Tags            []string  `json:"tags"`
}

type UpdateNotePayload struct {
	LinkedProjectID uuid.UUID `json:"linkedProjectID"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Favorited       string    `json:"favorited"`
	Tags            []string  `json:"tags"`
}
