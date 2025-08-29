package auth

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"

	"roombooker/internal/config"
	"roombooker/internal/repository"
)

type Service struct {
	repo    *repository.Repository
	config  *config.Config
	jwtAuth *jwtauth.JWTAuth
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo:    repo,
		config:  cfg,
		jwtAuth: jwtauth.New("HS256", []byte(cfg.Auth.JWTSecret), nil),
	}
}

func (s *Service) HashPassword(password string) []byte {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		// Fallback to fixed salt if random fails
		salt = []byte("some-random-salt-fixed-fallback")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hash
}

func (s *Service) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	_, tokenString, err := s.jwtAuth.Encode(claims)
	return tokenString, err
}

func (s *Service) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := s.jwtAuth.Decode(tokenString)
	if err != nil {
		return nil, err
	}

	claims, err := token.AsMap(context.Background())
	if err != nil {
		return nil, err
	}

	mapClaims := jwt.MapClaims(claims)
	return &mapClaims, nil
}
