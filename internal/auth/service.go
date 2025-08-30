package auth

import (
	"context"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return h
}

func (s *Service) VerifyPassword(hash []byte, password string) bool {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return false
	}
	return true
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
