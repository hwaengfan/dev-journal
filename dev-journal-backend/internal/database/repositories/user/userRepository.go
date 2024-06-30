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
	query := "SELECT * FROM users WHERE id = ?"
	row := store.database.QueryRow(query, id)

	user := &userModel.User{}
	error := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if error != nil {
		return nil, error
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (store *Store) GetUserByEmail(email string) (*userModel.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row := store.database.QueryRow(query, email)

	user := &userModel.User{}
	error := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if error != nil {
		return nil, error
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
