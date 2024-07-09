package taskRepository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	taskModel "github.com/hwaengfan/dev-journal-backend/internal/models/task"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{database: database}
}

// CreateTask creates a new task
func (store *Store) CreateTask(task taskModel.Task) (uuid.UUID, error) {
	taskID := uuid.New()

	query := "INSERT INTO tasks (id, linkedProjectID, description, completed) VALUES (?, ?, ?, ?)"
	_, error := store.database.Exec(query, taskID, task.LinkedProjectID, task.Description, task.Completed)
	if error != nil {
		return uuid.Nil, fmt.Errorf("failed to create task: %v", error)
	}

	return taskID, nil
}

// GetTasksByLinkedProjectID gets tasks by linked project ID
func (store *Store) GetTasksByLinkedProjectID(linkedProjectID uuid.UUID) ([]*taskModel.Task, error) {
	// query tasks by project ID
	query := "SELECT id, linkedProjectID, description, completed FROM tasks WHERE linkedProjectID = ?"
	rows, error := store.database.Query(query, linkedProjectID)
	if error != nil {
		return nil, fmt.Errorf("failed to get tasks by linked project ID: %v", error)
	}
	defer rows.Close()

	// scan tasks from rows
	tasks, error := scanTasksFromRows(rows)
	if error != nil {
		return nil, error
	}

	return tasks, nil
}

// GetTaskByID gets a task by its ID
func (store *Store) GetTaskByID(id uuid.UUID) (*taskModel.Task, error) {
	// query task by ID
	query := "SELECT id, linkedProjectID, description, completed FROM tasks WHERE id = ?"
	row := store.database.QueryRow(query, id)

	// scan task from row
	task, error := scanTaskFromRow(row)
	if error != nil {
		return nil, fmt.Errorf("failed to scan task from row: %v", error)
	}

	return task, nil
}

// UpdateTaskByID updates a task by its ID
func (store *Store) UpdateTaskByID(task taskModel.Task, id uuid.UUID) error {
	// base query
	query := "UPDATE tasks SET"
	var updates []string
	var args []interface{}

	// conditionally add fields to update
	if task.LinkedProjectID != uuid.Nil {
		updates = append(updates, "linkedProjectID = ?")
		args = append(args, task.LinkedProjectID)
	}
	if task.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, task.Description)
	}
	if task.Completed != "" {
		updates = append(updates, "completed = ?")
		args = append(args, task.Completed)
	}

	// check if there are fields to update
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// finalize query
	query += " " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	// execute the query
	_, err := store.database.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

// DeleteTaskByID deletes a task by its ID
func (store *Store) DeleteTaskByID(id uuid.UUID) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, error := store.database.Exec(query, id)
	if error != nil {
		return fmt.Errorf("failed to delete task: %v", error)
	}

	return nil
}

// DeleteTasksByLinkedProjectID deletes all tasks by linked project ID
func (store *Store) DeleteTasksByLinkedProjectID(linkedProjectID uuid.UUID) error {
	query := "DELETE FROM tasks WHERE linkedProjectID = ?"
	_, error := store.database.Exec(query, linkedProjectID)
	if error != nil {
		return fmt.Errorf("failed to delete all tasks by linked project ID: %v", error)
	}

	return nil
}

// scanTaskFromRows scans MySQL rows into a slice of task objects
func scanTasksFromRows(rows *sql.Rows) ([]*taskModel.Task, error) {
	tasks := make([]*taskModel.Task, 0)
	for rows.Next() {
		task := new(taskModel.Task)

		error := rows.Scan(&task.ID, &task.LinkedProjectID, &task.Description, &task.Completed)
		if error != nil {
			return nil, fmt.Errorf("failed to scan project from rows: %v", error)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// scanTask a MySQL row into a task object
func scanTaskFromRow(row *sql.Row) (*taskModel.Task, error) {
	task := new(taskModel.Task)

	error := row.Scan(&task.ID, &task.LinkedProjectID, &task.Description, &task.Completed)
	if error == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	} else if error != nil {
		return nil, fmt.Errorf("failed to scan project from row: %v", error)
	}

	return task, nil
}
