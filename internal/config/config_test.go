package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "sqlite3", cfg.Database.Driver)
	assert.Equal(t, "file:roombooker.db?cache=shared&_fk=1", cfg.Database.DSN)
	assert.Equal(t, "your-secret-key", cfg.Auth.JWTSecret)
	assert.Equal(t, "http://localhost:8080", cfg.App.BaseURL)
	assert.Equal(t, "America/New_York", cfg.App.OfficeTZ)
}

func TestLoad_Environment(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_DRIVER", "postgres")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	os.Setenv("JWT_SECRET", "custom-secret")
	os.Setenv("APP_BASE_URL", "https://example.com")
	os.Setenv("OFFICE_TZ", "Europe/London")
	defer os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, "postgres://user:pass@localhost/db", cfg.Database.DSN)
	assert.Equal(t, "custom-secret", cfg.Auth.JWTSecret)
	assert.Equal(t, "https://example.com", cfg.App.BaseURL)
	assert.Equal(t, "Europe/London", cfg.App.OfficeTZ)
}
