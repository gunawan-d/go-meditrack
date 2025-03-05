package repository

import (
	"meditrack/structs"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create user
func (r *UserRepository) CreateUser(user structs.User) error {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Role)
	return err
}

// Get by email
func (r *UserRepository) GetUserByemail(email string) (*structs.User, error) {
	user := structs.User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := r.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Check existing users
func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.DB.Get(&exists, query, email)
	return exists, err
}
