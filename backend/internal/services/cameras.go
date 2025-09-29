package services

import (
	"context"
	"errors"

	"smart-city-surveillance/internal/models"

	"gorm.io/gorm"
)

type CameraService interface {
	GetAll(ctx context.Context, userRole models.UserRole) ([]models.Camera, error)
	GetByID(ctx context.Context, id string,  userId string, userRole models.UserRole) (*models.Camera, error)
	GetByPremiseID(ctx context.Context, premiseID string, userRole models.UserRole) ([]models.Camera, error)
	GetAssignedByGuardID(ctx context.Context, guardID string) ([]models.Camera, error)
	UpdateStatus(ctx context.Context, id string, status models.CameraStatus, userRole models.UserRole) error
}

type cameraService struct {
	db *gorm.DB
}

func NewCameraService(db *gorm.DB) CameraService {
	return &cameraService{db: db}
}

// GetAll returns all cameras; only SCS Operator can access all cameras.
func (s *cameraService) GetAll(ctx context.Context, userRole models.UserRole) ([]models.Camera, error) {
	if userRole != models.RoleSCSOperator {
		return nil, errors.New("permission denied")
	}
	var cameras []models.Camera
	err := s.db.WithContext(ctx).Preload("Premise").Preload("Guards").Find(&cameras).Error
	return cameras, err
}

// GetByID returns camera details by ID; SCS Operator can access all; Security Guard can access assigned cameras only.
func (s *cameraService) GetByID(ctx context.Context, id string, userId string, userRole models.UserRole) (*models.Camera, error) {
	var camera *models.Camera
	if userRole == models.RoleSecurityGuard {
		// Guards can only access assigned cameras
		err := s.db.WithContext(ctx).
		Joins("JOIN camera_guards ON cameras.id = camera_guards.camera_id").
		Where("camera_guards.guard_id = ? and cameras.id = ?", userId, id).
		Preload("Premise").
		First(&camera).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.WithContext(ctx).Preload("Premise").Preload("Guards").First(&camera, "id = ?", id).Error
			if err != nil {
				return nil, err
			}
	}
	return camera, nil
}

// GetByPremiseID returns cameras for a premise; only SCS Operator can access
func (s *cameraService) GetByPremiseID(ctx context.Context, premiseID string, userRole models.UserRole) ([]models.Camera, error) {
	if userRole != models.RoleSCSOperator {
		return nil, errors.New("permission denied")
	}
	var cameras []models.Camera
	err := s.db.WithContext(ctx).Where("premise_id = ?", premiseID).Preload("Guards").Find(&cameras).Error
	return cameras, err
}

// GetAssignedByGuardID returns cameras assigned to a guard
func (s *cameraService) GetAssignedByGuardID(ctx context.Context, guardID string) ([]models.Camera, error) {
	var cameras []models.Camera
	err := s.db.WithContext(ctx).
		Joins("JOIN camera_guards ON cameras.id = camera_guards.camera_id").
		Where("camera_guards.guard_id = ?", guardID).
		Preload("Premise").
		Find(&cameras).Error
	return cameras, err
}

// UpdateStatus updates camera status; only SCS Operator can update status
func (s *cameraService) UpdateStatus(ctx context.Context, id string, status models.CameraStatus, userRole models.UserRole) error {
	if userRole != models.RoleSCSOperator {
		return errors.New("permission denied")
	}
	return s.db.WithContext(ctx).Model(&models.Camera{}).Where("id = ?", id).Update("status", status).Error
}
