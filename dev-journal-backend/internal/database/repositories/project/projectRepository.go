package projectRepository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	projectModel "github.com/hwaengfan/dev-journal-backend/internal/models/project"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{database: database}
}

// CreateProject creates a new project
func (store *Store) CreateProject(project projectModel.Project) (uuid.UUID, error) {
	projectID := uuid.New()

	query := "INSERT INTO projects (id, userID, title, description, priority, deadline) VALUES (?, ?, ?, ?, ?, ?)"
	_, error := store.database.Exec(query, projectID, project.UserID, project.Title, project.Description, project.Priority, project.Deadline)
	if error != nil {
		return uuid.Nil, fmt.Errorf("failed to create project: %v", error)
	}

	return projectID, nil
}

// GetProjectsByUserID retrieves all projects by a user's ID
func (store *Store) GetProjectsByUserID(userID uuid.UUID) ([]*projectModel.Project, error) {
	// query projects by user ID
	query := "SELECT id, title, description, priority, deadline, dateCreated, lastEdited FROM projects WHERE userID = ?"
	rows, error := store.database.Query(query, userID)
	if error != nil {
		return nil, fmt.Errorf("failed to get projects by user ID: %v", error)
	}
	defer rows.Close()

	// scan projects from rows
	projects, error := scanProjectFromRows(rows)
	if error != nil {
		return nil, error
	}

	return projects, nil
}

// GetProjectByID retrieves a project by its ID and user's ID
func (store *Store) GetProjectByID(id uuid.UUID, userID uuid.UUID) (*projectModel.Project, error) {
	// query project by ID
	query := "SELECT id, title, description, priority, deadline, dateCreated, lastEdited FROM projects WHERE id = ? AND userID = ?"
	row := store.database.QueryRow(query, id, userID)

	// scan project from row
	project, error := scanProjectFromRow(row)
	if error != nil {
		return nil, fmt.Errorf("failed to scan project from row: %v", error)
	}

	return project, nil
}

// UpdateProject updates a project by its ID and user's ID
func (store *Store) UpdateProject(project projectModel.Project, id uuid.UUID, userID uuid.UUID) error {
	// base query
	query := "UPDATE projects SET"
	var updates []string
	var args []interface{}

	// conditionally add fields to update
	if project.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, project.Title)
	}
	if project.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, project.Description)
	}
	if project.Priority != "" {
		updates = append(updates, "priority = ?")
		args = append(args, project.Priority)
	}
	if project.Deadline != "" {
		updates = append(updates, "deadline = ?")
		args = append(args, project.Deadline)
	}

	// check if there are fields to update
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// finalize query
	query += " " + strings.Join(updates, ", ") + " WHERE id = ? AND userID = ?"
	args = append(args, id, userID)

	// execute the query
	_, err := store.database.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	return nil
}

// DeleteProject deletes a project by its ID and user's ID
func (store *Store) DeleteProject(id uuid.UUID, userID uuid.UUID) error {
	query := "DELETE FROM projects WHERE id = ? AND userID = ?"
	_, error := store.database.Exec(query, id, userID)
	if error != nil {
		return fmt.Errorf("failed to delete project: %v", error)
	}

	return nil
}

// scanProjectFromRows scans MySQL rows into a slice of project objects
func scanProjectFromRows(rows *sql.Rows) ([]*projectModel.Project, error) {
	projects := make([]*projectModel.Project, 0)
	for rows.Next() {
		project := new(projectModel.Project)

		error := rows.Scan(&project.ID, &project.Title, &project.Description, &project.Priority, &project.Deadline, &project.DateCreated, &project.LastEdited)
		if error != nil {
			return nil, fmt.Errorf("failed to scan project from rows: %v", error)
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// scanProjectFromRow scans a MySQL row into a new project object
func scanProjectFromRow(row *sql.Row) (*projectModel.Project, error) {
	project := new(projectModel.Project)
	error := row.Scan(&project.ID, &project.Title, &project.Description, &project.Priority, &project.Deadline, &project.DateCreated, &project.LastEdited)

	if error == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	} else if error != nil {
		return nil, fmt.Errorf("failed to scan project from row: %v", error)
	}

	return project, nil
}
