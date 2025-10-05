package services

import (
	"context"
	"errors"
	"fmt"

	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/pkg/websocket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AlertsService defines alert-related operations
type AlertsService interface {
	GetAlerts(ctx context.Context, filters AlertsFilter, userRole models.UserRole, userID string) ([]models.Alert, error)
	GetAlert(ctx context.Context, id string, userRole models.UserRole, userID string) (*models.Alert, error)
	AcknowledgeAlert(ctx context.Context, id string, userRole models.UserRole) (*models.Alert, error)
	AssignAlert(ctx context.Context, id string, guardID []string) (*models.Alert, *models.Incident, error)
	CreateAlert(ctx context.Context, alert models.Alert) (*models.Alert, error)
	UpdateAlert(ctx context.Context, id string, status models.AlertStatus) (*models.Alert, error)
}

// AlertsFilter contains optional filter parameters for listing alerts
type AlertsFilter struct {
	Status    string
	Severity  string
	Type      string
	PremiseID string
}

type alertsService struct {
	db    *gorm.DB
	wsHub *websocket.Hub
}

func NewAlertsService(db *gorm.DB, wsHub *websocket.Hub) AlertsService {
	return &alertsService{db: db, wsHub: wsHub}
}

func (s *alertsService) GetAlerts(ctx context.Context, filters AlertsFilter, userRole models.UserRole, userID string) ([]models.Alert, error) {
	var alerts []models.Alert
	query := s.db.WithContext(ctx).Preload("Camera").Preload("Premise").Preload("AssignedGuard")

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Severity != "" {
		query = query.Where("severity = ?", filters.Severity)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.PremiseID != "" {
		if premiseUUID, err := uuid.Parse(filters.PremiseID); err == nil {
			query = query.Where("premise_id = ?", premiseUUID)
		}
	}

	if userRole == models.RoleSecurityGuard {
		query = query.Where("assigned_guard_id = ?", userID)
	}

	query = query.Order("created_at DESC")
	if err := query.Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}

func (s *alertsService) GetAlert(ctx context.Context, id string, userRole models.UserRole, userID string) (*models.Alert, error) {
	alertID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var alert models.Alert
	query := s.db.WithContext(ctx).Preload("Camera").Preload("Premise").Preload("AssignedGuard").Preload("Incident")
	if userRole == models.RoleSecurityGuard {
		query = query.Where("assigned_guard_id = ?", userID)
	}
	if err := query.First(&alert, "id = ?", alertID).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (s *alertsService) AcknowledgeAlert(ctx context.Context, id string, userRole models.UserRole) (*models.Alert, error) {
	if userRole != models.RoleSCSOperator {
		return nil, gorm.ErrInvalidData
	}
	alertID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var alert models.Alert
	if err := s.db.WithContext(ctx).First(&alert, "id = ?", alertID).Error; err != nil {
		return nil, err
	}
	alert.Status = models.AlertStatusAcknowledged
	if err := s.db.WithContext(ctx).Save(&alert).Error; err != nil {
		return nil, err
	}
	s.wsHub.Broadcast("alert_acknowledged", alert)
	return &alert, nil
}

func (s *alertsService) AssignAlert(ctx context.Context, id string, guardIDs []string) (*models.Alert, *models.Incident, error) {
	alertUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid alert id: %w", err)
	}

	// Lấy alert
	var alert models.Alert
	if err := s.db.WithContext(ctx).First(&alert, "id = ?", alertUUID).Error; err != nil {
		return nil, nil, fmt.Errorf("alert not found: %w", err)
	}

	// Validate guard IDs
	if len(guardIDs) == 0 {
		return nil, nil, errors.New("no guard IDs provided")
	}

	var guards []models.User
	if err := s.db.WithContext(ctx).
		Where("id IN ? AND role = ?", guardIDs, models.RoleSecurityGuard).
		Find(&guards).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to find guards: %w", err)
	}
	if len(guards) == 0 {
		return nil, nil, errors.New("no valid guards found")
	}

	// Tạo incident
	incident := models.Incident{
		AlertID:     alertUUID,
		Status:      models.IncidentStatusOpen,
		Location:    alert.Location,
		Description: alert.Description,
	}
	if err := s.db.WithContext(ctx).Create(&incident).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to create incident: %w", err)
	}

	// Cập nhật alert status
	alert.Status = models.AlertStatusAssigned
	if err := s.db.WithContext(ctx).Save(&alert).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to update alert: %w", err)
	}

	// Insert nhiều incident_guards
	var incidentGuards []models.IncidentGuard
	for _, g := range guards {
		incidentGuards = append(incidentGuards, models.IncidentGuard{
			IncidentID: incident.ID,
			GuardID:    g.ID,
		})
	}
	if err := s.db.WithContext(ctx).Create(&incidentGuards).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to assign guards: %w", err)
	}

	// Gửi notification
	for _, g := range guards {
		s.wsHub.SendToUser(g.ID.String(), "guard_dispatched", map[string]any{
			"alert_id":    alert.ID,
			"incident_id": incident.ID,
			"title":       alert.Title,
			"description": alert.Description,
			"location":    alert.Location,
			"severity":    alert.Severity,
		})
	}
	s.wsHub.BroadcastToRole("scs_operator", "alert_assigned", map[string]any{
		"alert_id":    alert.ID,
		"incident_id": incident.ID,
		"guards":      guards,
	})

	return &alert, &incident, nil
}


func (s *alertsService) CreateAlert(ctx context.Context, alert models.Alert) (*models.Alert, error) {
	alert.Status = models.AlertStatusPending
	if err := s.db.WithContext(ctx).Create(&alert).Error; err != nil {
		return nil, err
	}
	s.wsHub.BroadcastToRole("scs_operator", "alert_created", alert)
	return &alert, nil
}

func (s *alertsService) UpdateAlert(ctx context.Context, id string, status models.AlertStatus) (*models.Alert, error) {
	alertID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var alert models.Alert
	if err := s.db.WithContext(ctx).First(&alert, "id = ?", alertID).Error; err != nil {
		return nil, err
	}
	alert.Status = status
	if err := s.db.WithContext(ctx).Save(&alert).Error; err != nil {
		return nil, err
	}
	s.wsHub.Broadcast("alert_updated", alert)
	return &alert, nil
}
