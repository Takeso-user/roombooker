package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"roombooker/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestHandler_HealthCheck(t *testing.T) {
	cfg := &config.Config{}
	handler := NewHandler(nil, nil, nil, cfg, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
	assert.Contains(t, w.Body.String(), "ok")
}

func TestHandler_AuthMiddleware_NoToken(t *testing.T) {
	cfg := &config.Config{}
	handler := NewHandler(nil, nil, nil, cfg, nil)

	req := httptest.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()

	next := handler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	next.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_AuthMiddleware_WithToken(t *testing.T) {
	cfg := &config.Config{}
	handler := NewHandler(nil, nil, nil, cfg, nil)

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	w := httptest.NewRecorder()

	next := handler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	next.ServeHTTP(w, req)

	// Should proceed to next handler (though token validation would fail in real scenario)
	assert.Equal(t, http.StatusOK, w.Code)
}
