package taskModel

import "github.com/google/uuid"

type Task struct {
	ID              uuid.UUID `json:"id"`
	LinkedProjectID uuid.UUID `json:"linkedProjectID"`
	Description     string    `json:"description"`
	Completed       string    `json:"completed"`
}

type TaskStore interface {
	CreateTask(task Task) (uuid.UUID, error)
	GetTasksByLinkedProjectID(linkedProjectID uuid.UUID) ([]*Task, error)
	GetTaskByID(id uuid.UUID) (*Task, error)
	UpdateTaskByID(task Task, id uuid.UUID) error
	DeleteTaskByID(id uuid.UUID) error
	DeleteTasksByLinkedProjectID(linkedProjectID uuid.UUID) error
}

type CreateTaskPayload struct {
	LinkedProjectID uuid.UUID `json:"linkedProjectID" validate:"required"`
	Description     string    `json:"description" validate:"required"`
	Completed       string    `json:"completed"` // default is false so no need to require it
}

type UpdateTaskPayload struct {
	LinkedProjectID uuid.UUID `json:"linkedProjectID"`
	Description     string    `json:"description"`
	Completed       string    `json:"completed"`
}
