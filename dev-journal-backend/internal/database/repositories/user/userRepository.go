package userRepository

import (
	"database/sql"
	"fmt"

	"github.com/hwaengfan/dev-journal-backend/internal/models/user"
)

type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{database: database}
}

// CreateUser creates a new user
func (store *Store) CreateUser(user userModel.User) error {
	return nil
}

// GetUserByEmail retrieves a user by email
func (store *Store) GetUserByEmail(email string) (*userModel.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := store.database.QueryRow(query, email)

	user := &userModel.User{}
	error := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if error != nil {
		return nil, error
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (store *Store) GetUserByID(id int) (*userModel.User, error) {
	return nil, nil
}
