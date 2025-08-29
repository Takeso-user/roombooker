package repository

import (
	"database/sql"
)

type Repository struct {
	db     *sql.DB
	driver string
}

func New(db *sql.DB, driver string) *Repository {
	return &Repository{db: db, driver: driver}
}

func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) Driver() string {
	return r.driver
}

// Placeholder for user methods
type User struct {
	ID       string
	Email    string
	Role     string
	Timezone string
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	var user User
	query := "SELECT id, email, role, timezone FROM users WHERE id = $1"
	if r.driver == "sqlite3" {
		query = "SELECT id, email, role, timezone FROM users WHERE id = ?"
	}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Role, &user.Timezone)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Add more repository methods as needed
