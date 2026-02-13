package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret []byte) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("email already exists")
		}
		return nil, err
	}

	token, err := s.generateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  *createdUser,
	}, nil
}

func (s *AuthService) generateToken(userID int, email string) (string, error) {
	claims := models.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
