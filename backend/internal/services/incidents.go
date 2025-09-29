package services

import (
	"context"

	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/pkg/websocket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IncidentsService defines operations on incidents
type IncidentsService interface {
	GetIncidents(ctx context.Context, userRole models.UserRole, userID string, status string) ([]models.Incident, error)
	GetIncident(ctx context.Context, id string, userRole models.UserRole, userID string) (*models.Incident, error)
	UpdateIncident(ctx context.Context, id string, status models.IncidentStatus, userRole models.UserRole, userID string) (*models.Incident, error)
	AddIncidentUpdate(ctx context.Context, incidentID string, update models.IncidentUpdate, userRole models.UserRole, userID string) (*models.IncidentUpdate, error)
	GetIncidentByAlertID(ctx context.Context, alertID string) (*models.Incident, error)
}

type incidentsService struct {
	db    *gorm.DB
	wsHub *websocket.Hub
}

func NewIncidentsService(db *gorm.DB, wsHub *websocket.Hub) IncidentsService {
	return &incidentsService{db: db, wsHub: wsHub}
}

func (s *incidentsService) GetIncidents(ctx context.Context, userRole models.UserRole, userID string, status string) ([]models.Incident, error) {
	var incidents []models.Incident
	query := s.db.WithContext(ctx).Preload("Alert").Preload("AssignedGuard").Preload("Updates")
	if userRole == models.RoleSecurityGuard {
		query = query.Where("assigned_guard_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query = query.Order("created_at DESC")
	if err := query.Find(&incidents).Error; err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *incidentsService) GetIncident(ctx context.Context, id string, userRole models.UserRole, userID string) (*models.Incident, error) {
	incidentID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var incident models.Incident
	query := s.db.WithContext(ctx).Preload("Alert").Preload("AssignedGuard").Preload("Updates")
	if userRole == models.RoleSecurityGuard {
		query = query.Where("assigned_guard_id = ?", userID)
	}
	if err := query.First(&incident, "id = ?", incidentID).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}

func (s *incidentsService) UpdateIncident(ctx context.Context, id string, status models.IncidentStatus, userRole models.UserRole, userID string) (*models.Incident, error) {
	incidentID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var incident models.Incident
	if err := s.db.WithContext(ctx).First(&incident, "id = ?", incidentID).Error; err != nil {
		return nil, err
	}
	if userRole == models.RoleSecurityGuard && incident.AssignedGuardID.String() != userID {
		return nil, gorm.ErrInvalidData
	}
	incident.Status = status
	if err := s.db.WithContext(ctx).Save(&incident).Error; err != nil {
		return nil, err
	}
	s.wsHub.Broadcast("incident_updated", incident)
	return &incident, nil
}

func (s *incidentsService) AddIncidentUpdate(ctx context.Context, incidentID string, update models.IncidentUpdate, userRole models.UserRole, userID string) (*models.IncidentUpdate, error) {
	iid, err := uuid.Parse(incidentID)
	if err != nil {
		return nil, err
	}
	var incident models.Incident
	if err := s.db.WithContext(ctx).First(&incident, "id = ?", iid).Error; err != nil {
		return nil, err
	}
	if userRole == models.RoleSecurityGuard && incident.AssignedGuardID.String() != userID {
		return nil, gorm.ErrInvalidData
	}
	update.IncidentID = iid
	if guardUUID, err := uuid.Parse(userID); err == nil {
		update.GuardID = guardUUID
	}
	if err := s.db.WithContext(ctx).Create(&update).Error; err != nil {
		return nil, err
	}
	if update.Type == models.UpdateTypeResolution {
		incident.Status = models.IncidentStatusResolved
		_ = s.db.WithContext(ctx).Save(&incident).Error
	}
	s.wsHub.BroadcastToRole("scs_operator", "incident_update_received", map[string]any{
		"incident_id": iid,
		"update":      update,
		"guard_name":  userID,
	})
	return &update, nil
}

func (s *incidentsService) GetIncidentByAlertID(ctx context.Context, alertID string) (*models.Incident, error) {
	aid, err := uuid.Parse(alertID)
	if err != nil {
		return nil, err
	}
	var incident models.Incident
	if err := s.db.WithContext(ctx).Preload("Alert").Preload("AssignedGuard").Preload("Updates").First(&incident, "alert_id = ?", aid).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}
