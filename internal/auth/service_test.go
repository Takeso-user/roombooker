package auth

import (
	"testing"

	"roombooker/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestService_HashPassword(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}
	service := NewService(nil, cfg)

	hash1 := service.HashPassword("password123")
	hash2 := service.HashPassword("password123")

	// Argon2 should produce different hashes for same input due to salt
	assert.NotEqual(t, hash1, hash2)
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	assert.Len(t, hash1, 32) // Argon2 output length
	assert.Len(t, hash2, 32)
}

func TestService_GenerateToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}
	service := NewService(nil, cfg)

	token, err := service.GenerateToken("user-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_ValidateToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}
	service := NewService(nil, cfg)

	token, err := service.GenerateToken("user-123")
	assert.NoError(t, err)

	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-123", (*claims)["user_id"])
}

func TestService_ValidateToken_Invalid(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}
	service := NewService(nil, cfg)

	_, err := service.ValidateToken("invalid-token")
	assert.Error(t, err)
}
