package projectModel

import "github.com/google/uuid"

type Project struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Deadline    string    `json:"deadline"`
	DateCreated string    `json:"dateCreated"`
	LastEdited  string    `json:"lastEdited"`
}

type ProjectStore interface {
	CreateProject(project Project) error
	GetProjectsByUserID(userID uuid.UUID) ([]*Project, error)
	GetProjectByID(id uuid.UUID, userID uuid.UUID) (*Project, error)
	UpdateProject(project Project, id uuid.UUID, userID uuid.UUID) error
	DeleteProject(id uuid.UUID, userID uuid.UUID) error
}

type CreateProjectPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	Deadline    string `json:"deadline" validate:"required"`
}

type UpdateProjectPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Deadline    string `json:"deadline"`
}
