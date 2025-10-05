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
	query := s.db.WithContext(ctx).
		Model(&models.Incident{}).
		Preload("Alert").
		Preload("AssignedGuards").
		Preload("Updates")

	if userRole == models.RoleSecurityGuard {
		// join bảng trung gian incident_guards để lọc các incident có guard tương ứng
		query = query.Joins("JOIN incident_guards ig ON ig.incident_id = incidents.id").
			Where("ig.guard_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("incidents.created_at DESC")

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

	query := s.db.WithContext(ctx).
		Model(&models.Incident{}).
		Preload("Alert").
		Preload("AssignedGuards").
		Preload("Updates").
		Where("incidents.id = ?", incidentID)

	if userRole == models.RoleSecurityGuard {
		query = query.Joins("JOIN incident_guards ig ON ig.incident_id = incidents.id").
			Where("ig.guard_id = ?", userID)
	}

	if err := query.First(&incident).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}
func (s *incidentsService) UpdateIncident(
	ctx context.Context,
	id string,
	status models.IncidentStatus,
	userRole models.UserRole,
	userID string,
) (*models.Incident, error) {
	incidentID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var incident models.Incident
	if err := s.db.WithContext(ctx).First(&incident, "id = ?", incidentID).Error; err != nil {
		return nil, err
	}

	// ✅ Nếu là guard, kiểm tra có thuộc incident_guards không
	if userRole == models.RoleSecurityGuard {
		var count int64
		if err := s.db.WithContext(ctx).
			Table("incident_guards").
			Where("incident_id = ? AND guard_id = ?", incidentID, userID).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, gorm.ErrInvalidData // không được phép
		}
	}

	incident.Status = status
	if err := s.db.WithContext(ctx).Save(&incident).Error; err != nil {
		return nil, err
	}

	s.wsHub.Broadcast("incident_updated", incident)
	return &incident, nil
}


func (s *incidentsService) AddIncidentUpdate(
	ctx context.Context,
	incidentID string,
	update models.IncidentUpdate,
	userRole models.UserRole,
	userID string,
) (*models.IncidentUpdate, error) {
	iid, err := uuid.Parse(incidentID)
	if err != nil {
		return nil, err
	}

	var incident models.Incident
	if err := s.db.WithContext(ctx).First(&incident, "id = ?", iid).Error; err != nil {
		return nil, err
	}

	// ✅ Nếu là guard, kiểm tra có trong incident_guards không
	if userRole == models.RoleSecurityGuard {
		var count int64
		if err := s.db.WithContext(ctx).
			Table("incident_guards").
			Where("incident_id = ? AND guard_id = ?", iid, userID).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, gorm.ErrInvalidData
		}
	}

	update.IncidentID = iid
	if guardUUID, err := uuid.Parse(userID); err == nil {
		update.GuardID = guardUUID
	}

	if err := s.db.WithContext(ctx).Create(&update).Error; err != nil {
		return nil, err
	}

	// ✅ Nếu update là loại resolution thì đổi status incident
	if update.Type == models.UpdateTypeResolution {
		incident.Status = models.IncidentStatusResolved
		if err := s.db.WithContext(ctx).Save(&incident).Error; err != nil {
			return nil, err
		}
	}

	// ✅ Broadcast event
	s.wsHub.BroadcastToRole("scs_operator", "incident_update_received", map[string]any{
		"incident_id": iid,
		"update":      update,
		"guard_id":    userID,
	})

	return &update, nil
}


func (s *incidentsService) GetIncidentByAlertID(ctx context.Context, alertID string) (*models.Incident, error) {
	aid, err := uuid.Parse(alertID)
	if err != nil {
		return nil, err
	}
	var incident models.Incident
	if err := s.db.WithContext(ctx).
		Preload("Alert").
		Preload("AssignedGuards").
		Preload("Updates").
		First(&incident, "alert_id = ?", aid).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}