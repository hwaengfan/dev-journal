package userRepository

import (
	"database/sql"
	"fmt"

	userModel "github.com/hwaengfan/dev-journal-backend/internal/models/user"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{database: database}
}

// CreateUser creates a new user
func (store *Store) CreateUser(user userModel.User) error {
	query := "INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)"
	_, error := store.database.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if error != nil {
		return error
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (store *Store) GetUserByID(id int) (*userModel.User, error) {
	// query user by ID
	query := "SELECT * FROM users WHERE id = ?"
	row, error := store.database.Query(query, id)
	if error != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", error)
	}
	defer row.Close()

	// scan user from row
	user, error := scanUserFromRow(row)
	if error != nil {
		return nil, fmt.Errorf("failed to scan user from row: %v", error)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (store *Store) GetUserByEmail(email string) (*userModel.User, error) {
	// query user by email
	query := "SELECT * FROM users WHERE email = ?"
	row, error := store.database.Query(query, email)
	if error != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", error)
	}
	defer row.Close()

	// scan user from row
	user, error := scanUserFromRow(row)
	if error != nil {
		return nil, fmt.Errorf("failed to scan user from row: %v", error)
	}

	return user, nil
}

// scanUserFromRow scans a MySQL row into a new user object
func scanUserFromRow(row *sql.Rows) (*userModel.User, error) {
	user := new(userModel.User)
	error := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if error != nil {
		return nil, error
	}

	return user, nil
}
