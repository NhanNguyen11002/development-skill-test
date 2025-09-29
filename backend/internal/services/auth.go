package services

import (
	"context"

	"smart-city-surveillance/internal/config"
	"smart-city-surveillance/internal/middleware"
	"smart-city-surveillance/internal/models"

	"gorm.io/gorm"
)

// AuthService defines authentication-related operations
type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, models.User, error)
	GetUserByID(ctx context.Context, userID string) (models.User, error)
}

type authService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) AuthService {
	return &authService{db: db, config: cfg}
}

func (s *authService) Login(ctx context.Context, username string, password string) (string, models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil {
		return "", models.User{}, err
	}
	if !middleware.CheckPassword(password, user.Password) {
		return "", models.User{}, gorm.ErrInvalidData
	}
	token, err := middleware.GenerateToken(&user, s.config)
	if err != nil {
		return "", models.User{}, err
	}
	return token, user, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
