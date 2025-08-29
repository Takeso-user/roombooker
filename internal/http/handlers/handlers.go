package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"roombooker/internal/auth"
	"roombooker/internal/config"
	"roombooker/internal/msgraph"
	"roombooker/internal/repository"
)

type Handler struct {
	repo        *repository.Repository
	authService *auth.Service
	graphClient *msgraph.Client
	config      *config.Config
	logger      *zap.Logger
}

func NewHandler(repo *repository.Repository, authService *auth.Service, graphClient *msgraph.Client, cfg *config.Config, logger *zap.Logger) *Handler {
	return &Handler{
		repo:        repo,
		authService: authService,
		graphClient: graphClient,
		config:      cfg,
		logger:      logger,
	}
}

func SetupRoutes(r *chi.Mux, repo *repository.Repository, authService *auth.Service, graphClient *msgraph.Client, cfg *config.Config, logger *zap.Logger) {
	h := NewHandler(repo, authService, graphClient, cfg, logger)

	r.Get("/health", h.HealthCheck)

	// Auth routes
	r.Post("/auth/login", h.Login)
	r.Get("/auth/oidc/start", h.OIDCStart)
	r.Get("/auth/oidc/callback", h.OIDCCallback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/me", h.GetMe)
		r.Get("/rooms", h.GetRooms)
		r.Get("/rooms/{id}/calendar", h.GetRoomCalendar)
		r.Post("/rooms/{id}/bookings", h.CreateBooking)
		r.Get("/bookings/{id}", h.GetBooking)
		r.Patch("/bookings/{id}", h.UpdateBooking)
		r.Delete("/bookings/{id}", h.DeleteBooking)
		r.Get("/availability", h.GetAvailability)
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Validate token logic here
		next.ServeHTTP(w, r)
	})
}

// Placeholder handlers
func (h *Handler) Login(w http.ResponseWriter, r *http.Request)           { /* implement */ }
func (h *Handler) OIDCStart(w http.ResponseWriter, r *http.Request)       { /* implement */ }
func (h *Handler) OIDCCallback(w http.ResponseWriter, r *http.Request)    { /* implement */ }
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request)           { /* implement */ }
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request)        { /* implement */ }
func (h *Handler) GetRoomCalendar(w http.ResponseWriter, r *http.Request) { /* implement */ }
func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request)   { /* implement */ }
func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request)      { /* implement */ }
func (h *Handler) UpdateBooking(w http.ResponseWriter, r *http.Request)   { /* implement */ }
func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request)   { /* implement */ }
func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) { /* implement */ }
