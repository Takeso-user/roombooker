package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetUserByID(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Create test table
	_, err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		display_name TEXT,
		role TEXT NOT NULL DEFAULT 'user',
		timezone TEXT DEFAULT 'UTC'
	)`)
	assert.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`INSERT INTO users (id, email, display_name, role, timezone) 
		VALUES (?, ?, ?, ?, ?)`, "test-id", "test@example.com", "Test User", "user", "UTC")
	assert.NoError(t, err)

	repo := New(db, "sqlite3")

	user, err := repo.GetUserByID("test-id")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test-id", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "user", user.Role)
}

func TestRepository_GetUserByID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		display_name TEXT,
		role TEXT NOT NULL DEFAULT 'user',
		timezone TEXT DEFAULT 'UTC'
	)`)
	assert.NoError(t, err)

	repo := New(db, "sqlite3")

	user, err := repo.GetUserByID("non-existent")
	assert.Error(t, err)
	assert.Nil(t, user)
}
