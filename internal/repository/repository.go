package repository

import (
	"database/sql"
	"fmt"
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

// CreateUser inserts a new user and returns the new ID
func (r *Repository) CreateUser(email, displayName, role, passwordHash string) (string, error) {
	var query string
	if r.driver == "sqlite3" {
		query = "INSERT INTO users(email, display_name, role, password_hash) VALUES (?, ?, ?, ?)"
	} else {
		query = "INSERT INTO users(email, display_name, role, password_hash) VALUES ($1, $2, $3, $4)"
	}
	if _, err := r.db.Exec(query, email, displayName, role, passwordHash); err != nil {
		return "", err
	}
	// Try to select id by email
	sel := "SELECT id FROM users WHERE email = $1"
	if r.driver == "sqlite3" {
		sel = "SELECT id FROM users WHERE email = ?"
	}
	var id string
	if err := r.db.QueryRow(sel, email).Scan(&id); err != nil {
		return email, nil
	}
	return id, nil
}

// GetUserByEmail fetches a user by email
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	query := "SELECT id, email, role, timezone, display_name, password_hash FROM users WHERE email = $1"
	if r.driver == "sqlite3" {
		query = "SELECT id, email, role, timezone, display_name, password_hash FROM users WHERE email = ?"
	}
	var displayName sql.NullString
	var passwordHash sql.NullString
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Role, &user.Timezone, &displayName, &passwordHash)
	if err != nil {
		return nil, err
	}
	// Note: password_hash and display_name not stored in struct currently; user of this method can query DB directly if needed
	return &user, nil
}

// GetUserCredentials returns id, password_hash, role, display_name for an email
func (r *Repository) GetUserCredentials(email string) (id string, passwordHash string, role string, displayName string, err error) {
	query := "SELECT id, password_hash, role, display_name FROM users WHERE email = $1"
	if r.driver == "sqlite3" {
		query = "SELECT id, password_hash, role, display_name FROM users WHERE email = ?"
	}
	var ph sql.NullString
	var dn sql.NullString
	err = r.db.QueryRow(query, email).Scan(&id, &ph, &role, &dn)
	if err != nil {
		return "", "", "", "", err
	}
	if ph.Valid {
		passwordHash = ph.String
	}
	if dn.Valid {
		displayName = dn.String
	}
	return id, passwordHash, role, displayName, nil
}

// ListUsers returns a simple list of users
func (r *Repository) ListUsers() ([]User, error) {
	query := "SELECT id, email, role, timezone FROM users"
	if r.driver == "sqlite3" {
		query = "SELECT id, email, role, timezone FROM users"
	}
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.Timezone); err != nil {
			continue
		}
		out = append(out, u)
	}
	return out, nil
}

// UpdateUserRole updates a user's role
func (r *Repository) UpdateUserRole(id, role string) error {
	query := "UPDATE users SET role = $1 WHERE id = $2"
	if r.driver == "sqlite3" {
		query = "UPDATE users SET role = ? WHERE id = ?"
	}
	_, err := r.db.Exec(query, role, id)
	return err
}

// CreateOffice inserts a new office and returns its id (or name fallback)
func (r *Repository) CreateOffice(name, timezone string) (string, error) {
	query := "INSERT INTO offices(name, timezone) VALUES ($1, $2)"
	if r.driver == "sqlite3" {
		query = "INSERT INTO offices(name, timezone) VALUES (?, ?)"
	}
	res, err := r.db.Exec(query, name, timezone)
	if err != nil {
		return "", err
	}
	if id, err := res.LastInsertId(); err == nil && id > 0 {
		return fmt.Sprintf("%d", id), nil
	}
	return name, nil
}

// CreateFloor inserts a new floor under an office
func (r *Repository) CreateFloor(officeID string, number int, label string) (string, error) {
	query := "INSERT INTO floors(office_id, number, label) VALUES ($1, $2, $3)"
	if r.driver == "sqlite3" {
		query = "INSERT INTO floors(office_id, number, label) VALUES (?, ?, ?)"
	}
	res, err := r.db.Exec(query, officeID, number, label)
	if err != nil {
		return "", err
	}
	if id, err := res.LastInsertId(); err == nil && id > 0 {
		return fmt.Sprintf("%d", id), nil
	}
	return fmt.Sprintf("%s-%d", officeID, number), nil
}

// CreateRoom inserts a new room under a floor
func (r *Repository) CreateRoom(floorID, name string, capacity int, equipment string) (string, error) {
	query := "INSERT INTO rooms(floor_id, name, capacity, equipment) VALUES ($1, $2, $3, $4)"
	if r.driver == "sqlite3" {
		query = "INSERT INTO rooms(floor_id, name, capacity, equipment) VALUES (?, ?, ?, ?)"
	}
	res, err := r.db.Exec(query, floorID, name, capacity, equipment)
	if err != nil {
		return "", err
	}
	if id, err := res.LastInsertId(); err == nil && id > 0 {
		return fmt.Sprintf("%d", id), nil
	}
	return name, nil
}

// Add more repository methods as needed
