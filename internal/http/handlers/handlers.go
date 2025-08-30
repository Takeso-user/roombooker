package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	// in-memory bookings store for dev/testing
	bookings   map[string][]Booking
	bookingsMu sync.Mutex
}

// Booking represents a calendar booking returned to the frontend
type Booking struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
	Color  string `json:"color,omitempty"`
}

// CreateBookingRequest is the expected payload from the frontend
type CreateBookingRequest struct {
	Title     string   `json:"title"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Attendees []string `json:"attendees"`
	RoomID    string   `json:"room_id"`
}

func NewHandler(repo *repository.Repository, authService *auth.Service, graphClient *msgraph.Client, cfg *config.Config, logger *zap.Logger) *Handler {
	return &Handler{
		repo:        repo,
		authService: authService,
		graphClient: graphClient,
		config:      cfg,
		logger:      logger,
		bookings:    make(map[string][]Booking),
	}
}

func SetupRoutes(r *chi.Mux, repo *repository.Repository, authService *auth.Service, graphClient *msgraph.Client, cfg *config.Config, logger *zap.Logger) {
	h := NewHandler(repo, authService, graphClient, cfg, logger)

	r.Get("/health", h.HealthCheck)

	// Auth routes
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Get("/auth/oidc/start", h.OIDCStart)
	r.Get("/auth/oidc/callback", h.OIDCCallback)
	r.Post("/auth/logout", h.Logout)

	// Public routes
	r.Get("/login", h.LoginPage)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/me", h.GetMe)

		// API routes
		r.Route("/api", func(r chi.Router) {
			r.Get("/offices", h.GetOffices)
			r.Get("/offices/{officeId}/rooms", h.GetRoomsByOffice)
			r.Get("/rooms/{id}/bookings", h.GetRoomBookings)
			r.Post("/bookings", h.CreateBooking)
			r.Get("/bookings/{id}", h.GetBooking)
			r.Patch("/bookings/{id}", h.UpdateBooking)
			r.Delete("/bookings/{id}", h.DeleteBooking)

			// Admin routes
			r.Route("/admin", func(r chi.Router) {
				r.Use(h.AdminMiddleware)
				r.Get("/users", h.GetUsers)
				r.Patch("/users/{id}/role", h.UpdateUserRole)
				r.Post("/rooms", h.CreateRoom)
				r.Patch("/rooms/{id}", h.UpdateRoom)
				r.Delete("/rooms/{id}", h.DeleteRoom)
				r.Post("/offices", h.CreateOffice)
				r.Patch("/offices/{id}", h.UpdateOffice)
				r.Delete("/offices/{id}", h.DeleteOffice)
			})
		})
	})

	// Serve static files
	r.Get("/static/*", h.ServeStatic)

	// Main page
	r.Get("/", h.MainPage)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			// Try to get token from cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Validate token
		claims, err := h.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store user ID in context for later use
		ctx := context.WithValue(r.Context(), "user_id", (*claims)["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Placeholder handlers
func (h *Handler) OIDCStart(w http.ResponseWriter, r *http.Request) {
	// For testing purposes, redirect directly to callback with mock code
	mockCode := "test-auth-code"
	mockState := "test-state"

	callbackURL := fmt.Sprintf("%s?code=%s&state=%s", h.config.Auth.OIDCRedirectURL, mockCode, mockState)
	http.Redirect(w, r, callbackURL, http.StatusFound)
}

func (h *Handler) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	_ = r.URL.Query().Get("state") // Validate state in production

	if code == "" {
		http.Error(w, "No authorization code", http.StatusBadRequest)
		return
	}

	// In a real implementation, exchange code for tokens
	// For now, just create a JWT token
	token, err := h.authService.GenerateToken("test-user-id")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to main page
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// For now, return mock user data
	user := map[string]interface{}{
		"id":    userID,
		"email": "test@example.com",
		"name":  "Test User",
		"role":  "user",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Register allows creating a new local user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	pwHash := h.authService.HashPassword(req.Password)
	id, err := h.repo.CreateUser(req.Email, req.DisplayName, "user", string(pwHash))
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Auto-login after registration
	token, err := h.authService.GenerateToken(id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "email": req.Email})
}

// Login authenticates a local user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	id, pwHashStr, role, _, err := h.repo.GetUserCredentials(req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if pwHashStr == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if !h.authService.VerifyPassword([]byte(pwHashStr), req.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.authService.GenerateToken(id)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	json.NewEncoder(w).Encode(map[string]string{"id": id, "role": role})
}

// Logout handler
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// Login page handler
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/test.html")
}

// Main page handler
func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/index.html")
}

// Serve static files
func (h *Handler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/")))
	fs.ServeHTTP(w, r)
}

// Get offices
func (h *Handler) GetOffices(w http.ResponseWriter, r *http.Request) {
	// Mock data for now
	offices := []map[string]interface{}{
		{"id": "1", "name": "Main Office"},
		{"id": "2", "name": "Branch Office"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offices)
}

// Get rooms by office
func (h *Handler) GetRoomsByOffice(w http.ResponseWriter, r *http.Request) {
	officeID := chi.URLParam(r, "officeId")

	// Mock data for now
	rooms := []map[string]interface{}{
		{"id": "1", "name": "Conference Room A", "capacity": 10, "office_id": officeID},
		{"id": "2", "name": "Meeting Room B", "capacity": 6, "office_id": officeID},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

// Get room bookings
func (h *Handler) GetRoomBookings(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	// Try to parse optional from/to query params for filtering
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	var fromT, toT time.Time
	var err error
	if fromStr != "" {
		fromT, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			// ignore parse error and return all
			fromT = time.Time{}
		}
	}
	if toStr != "" {
		toT, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			toT = time.Time{}
		}
	}

	h.bookingsMu.Lock()
	defer h.bookingsMu.Unlock()

	raw := h.bookings[roomID]
	// If no bookings exist in store, return a small sample so calendar isn't empty
	if len(raw) == 0 {
		sample := Booking{
			ID:     "sample-1",
			Title:  "Team Meeting",
			Start:  "2024-01-15T10:00:00Z",
			End:    "2024-01-15T11:00:00Z",
			UserID: "user1",
			RoomID: roomID,
			Color:  "#3788d8",
		}
		raw = []Booking{sample}
	}

	// Filter by range if provided
	var out []Booking
	for _, b := range raw {
		if fromT.IsZero() && toT.IsZero() {
			out = append(out, b)
			continue
		}
		// parse booking times
		sb, err1 := time.Parse(time.RFC3339, b.Start)
		eb, err2 := time.Parse(time.RFC3339, b.End)
		if err1 != nil || err2 != nil {
			// skip invalid
			continue
		}
		// overlap check: b.Start < to && b.End > from
		if (toT.IsZero() || sb.Before(toT)) && (fromT.IsZero() || eb.After(fromT)) {
			out = append(out, b)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// Create booking
func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Build booking and store in-memory
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	userID := "unknown"
	if v := r.Context().Value("user_id"); v != nil {
		userID = fmt.Sprintf("%v", v)
	}
	// Normalize start/end times: accept RFC3339, datetime-local (no zone), or common variants and store as RFC3339
	var startT, endT time.Time
	var err error
	parseCandidates := []string{
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	tryParse := func(val string) (time.Time, error) {
		for _, layout := range parseCandidates {
			if t, e := time.Parse(layout, val); e == nil {
				return t, nil
			}
		}
		// If value looks like 'YYYY-MM-DDTHH:MM' without timezone, try appending Z
		if len(val) == 16 && val[10] == 'T' {
			if t, e := time.Parse(time.RFC3339, val+":00Z"); e == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("no parse")
	}

	startT, err = tryParse(req.StartTime)
	if err == nil {
		// if parsed without zone, assume local and convert to UTC
		if startT.Location() == time.Local {
			startT = startT.UTC()
		}
	}
	endT, err = tryParse(req.EndTime)
	if err == nil {
		if endT.Location() == time.Local {
			endT = endT.UTC()
		}
	}

	startStr := req.StartTime
	endStr := req.EndTime
	if !startT.IsZero() {
		startStr = startT.Format(time.RFC3339)
	}
	if !endT.IsZero() {
		endStr = endT.Format(time.RFC3339)
	}

	b := Booking{
		ID:     id,
		Title:  req.Title,
		Start:  startStr,
		End:    endStr,
		UserID: userID,
		RoomID: req.RoomID,
		Color:  "#3788d8",
	}

	h.bookingsMu.Lock()
	h.bookings[req.RoomID] = append(h.bookings[req.RoomID], b)
	h.bookingsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

// Admin middleware
func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id")
		if userID == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user is admin (mock check for now)
		// In real implementation, check user role from database
		isAdmin := true // Mock admin check

		if !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Admin handlers
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []map[string]interface{}{
		{"id": "1", "email": "user1@example.com", "name": "User One", "role": "user"},
		{"id": "2", "email": "admin@example.com", "name": "Admin User", "role": "admin"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Mock response
	response := map[string]interface{}{
		"id":   userID,
		"role": req["role"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FloorID   string `json:"floor_id"`
		Name      string `json:"name"`
		Capacity  int    `json:"capacity"`
		Equipment string `json:"equipment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	id, err := h.repo.CreateRoom(req.FloorID, req.Name, req.Capacity, req.Equipment)
	if err != nil {
		http.Error(w, "Failed to create room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "name": req.Name})
}

func (h *Handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	var room map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	room["id"] = roomID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Room " + roomID + " deleted"})
}

func (h *Handler) CreateOffice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Timezone string `json:"timezone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	id, err := h.repo.CreateOffice(req.Name, req.Timezone)
	if err != nil {
		http.Error(w, "Failed to create office: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "name": req.Name})
}

func (h *Handler) UpdateOffice(w http.ResponseWriter, r *http.Request) {
	officeID := chi.URLParam(r, "id")

	var office map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&office); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	office["id"] = officeID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(office)
}

func (h *Handler) DeleteOffice(w http.ResponseWriter, r *http.Request) {
	officeID := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Office " + officeID + " deleted"})
}

// Placeholder handlers (keeping for compatibility)
// func (h *Handler) Login(w http.ResponseWriter, r *http.Request)           { /* implemented above */ }
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request)        { /* implement */ }
func (h *Handler) GetRoomCalendar(w http.ResponseWriter, r *http.Request) { /* implement */ }
func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request)      { /* implement */ }
func (h *Handler) UpdateBooking(w http.ResponseWriter, r *http.Request)   { /* implement */ }
func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request)   { /* implement */ }
func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) { /* implement */ }
