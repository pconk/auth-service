package service

import (
	"auth-service/internal/config"
	"auth-service/internal/entity"
	"auth-service/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password, email, role, position string) error
	Login(username, password string) (string, error) // Returns JWT Token
	ValidateToken(tokenString string) (*CustomClaims, error)
	GetUserByID(id int) (*entity.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"` // Menambahkan username
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(cfg.JWT_SECRET),
	}
}

func (s *authService) Register(username, password, email, role, position string) error {
	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Role:         role,
		Position:     position,
	}

	return s.repo.Create(user)
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare Hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username, // Menambahkan username
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 Hari
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) GetUserByID(id int) (*entity.User, error) {
	return s.repo.FindByID(id)
}
