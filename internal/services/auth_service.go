package services

import (
	"errors"
	"bankapi/internal/config"
	"bankapi/internal/models"
	"bankapi/internal/repositories"
	"bankapi/pkg/crypto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUserExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

func RegisterUser(req models.RegisterRequest) (*models.User, error) {
	existing, _ := repositories.GetUserByUsername(req.Username)
	if existing != nil {
		return nil, ErrUserExists
	}
	existing, _ = repositories.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = repositories.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AuthenticateUser(req models.LoginRequest, cfg *config.Config) (string, error) {
	user, err := repositories.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		return "", ErrInvalidCredentials
	}

	if !crypto.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
