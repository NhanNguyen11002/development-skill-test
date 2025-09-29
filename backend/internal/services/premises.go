package services

import (
	"context"
	"fmt"
	"smart-city-surveillance/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PremisesService interface {
	GetPremises(ctx context.Context) ([]models.Premise, error)
	GetPremise(ctx context.Context, id uuid.UUID) (*models.Premise, error)
	GetPremiseCameras(ctx context.Context, id uuid.UUID) ([]models.Camera, error)
}

type premisesService struct {
	db *gorm.DB
}

func NewPremisesService(db *gorm.DB) PremisesService {
	return &premisesService{db: db}
}

func (s *premisesService) GetPremises(ctx context.Context) ([]models.Premise, error) {
	var premises []models.Premise
	if err := s.db.WithContext(ctx).Find(&premises).Error; err != nil {
		return nil, err
	}
	return premises, nil
}

func (s *premisesService) GetPremise(ctx context.Context, id uuid.UUID) (*models.Premise, error) {
	fmt.Printf("Id nè %s", id)
	var premise models.Premise
	if err := s.db.WithContext(ctx).First(&premise, id).Error; err != nil {
		return nil, err
	}
	return &premise, nil
}

func (s *premisesService) GetPremiseCameras(ctx context.Context, id uuid.UUID) ([]models.Camera, error) {
	var cameras []models.Camera
	if err := s.db.WithContext(ctx).
		Where("premise_id = ?", id).
		Find(&cameras).Error; err != nil {
		return nil, err
	}
	return cameras, nil
}
