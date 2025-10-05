package services

import (
	"context"
	"errors"

	"smart-city-surveillance/internal/models"

	"gorm.io/gorm"
)

type UserService interface {
	GetAll(ctx context.Context, userRole models.UserRole) ([]models.User, error)
	GetByAssignedCameraID(ctx context.Context, cameraID string, userRole models.UserRole) ([]models.User, error)
	GetByAssignedIncidentID(ctx context.Context, incidentID string, userRole models.UserRole) ([]models.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetAll(ctx context.Context, userRole models.UserRole) ([]models.User, error) {
	if userRole != models.RoleSCSOperator {
		return nil, errors.New("permission denied")
	}
	var users []models.User
	err := s.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (s *userService) GetByAssignedCameraID(ctx context.Context, cameraID string, userRole models.UserRole) ([]models.User, error) {
	if userRole != models.RoleSCSOperator {
		return nil, errors.New("permission denied")
	}
	var users []models.User
	err := s.db.WithContext(ctx).
		Joins("JOIN camera_guards ON users.id = camera_guards.guard_id").
		Where("camera_guards.camera_id = ?", cameraID).
		Find(&users).Error
	return users, err
}

func (s *userService) GetByAssignedIncidentID(ctx context.Context, incidentID string, userRole models.UserRole) ([]models.User, error) {
	if userRole != models.RoleSCSOperator {
		return nil, errors.New("permission denied")
	}

	var users []models.User
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Joins("JOIN incident_guards ig ON ig.guard_id = users.id").
		Where("ig.incident_id = ?", incidentID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
