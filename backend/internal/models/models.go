package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// =======================
// User & Role
// =======================

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	CameraAssignments   []CameraGuard   `json:"camera_assignments,omitempty" gorm:"foreignKey:GuardID"`
	IncidentAssignments []IncidentGuard `json:"incident_assignments,omitempty" gorm:"foreignKey:GuardID"`
}

type UserRole string

const (
	RoleSCSOperator   UserRole = "scs_operator"
	RoleSecurityGuard UserRole = "security_guard"
)

// =======================
// Premise & Camera
// =======================

type Premise struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string      `json:"name" gorm:"not null"`
	Address     string      `json:"address" gorm:"not null"`
	Type        PremiseType `json:"type" gorm:"not null"`
	FloorPlans  string      `json:"floor_plans"`
	Description string      `json:"description"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// Relationships
	Cameras []Camera `json:"cameras,omitempty" gorm:"foreignKey:PremiseID;references:ID"`
	Alerts  []Alert  `json:"alerts,omitempty" gorm:"foreignKey:PremiseID;references:ID"`
}

type PremiseType string

const (
	PremiseTypeOffice     PremiseType = "office"
	PremiseTypeSubstation PremiseType = "substation"
)

type Camera struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name           string    `json:"name" gorm:"not null"`
	Location       string    `json:"location" gorm:"not null"`
	StreamURL      string    `json:"stream_url" gorm:"not null"`
	Status         CameraStatus `json:"status" gorm:"default:'active'"`
	PremiseID      uuid.UUID    `json:"premise_id" gorm:"type:uuid;not null"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`

	// Relationships
	Premise Premise `json:"premise,omitempty" gorm:"foreignKey:PremiseID;references:ID"`
	Guards  []User  `json:"guards,omitempty" gorm:"many2many:camera_guards;joinForeignKey:CameraID;JoinReferences:GuardID"`
}

type CameraStatus string

const (
	CameraStatusActive       CameraStatus = "active"
	CameraStatusInactive     CameraStatus = "inactive"
	CameraStatusMaintenance  CameraStatus = "maintenance"
)

// =======================
// Alert & Incident
// =======================

type Alert struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type        AlertType     `json:"type" gorm:"not null"`
	Severity    AlertSeverity `json:"severity" gorm:"not null"`
	Title       string        `json:"title" gorm:"not null"`
	Description string        `json:"description" gorm:"not null"`
	Location    string        `json:"location" gorm:"not null"`
	Status      AlertStatus   `json:"status" gorm:"default:'pending'"`
	CameraID    *uuid.UUID    `json:"camera_id,omitempty" gorm:"type:uuid"`
	PremiseID   uuid.UUID     `json:"premise_id" gorm:"type:uuid;not null"`
	AssignedGuardID *uuid.UUID `json:"assigned_guard_id,omitempty" gorm:"type:uuid"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`

	// Relationships
	Camera   *Camera   `json:"camera,omitempty" gorm:"foreignKey:CameraID;references:ID"`
	Premise  Premise   `json:"premise,omitempty" gorm:"foreignKey:PremiseID;references:ID"`
	AssignedGuard *User `json:"assigned_guard,omitempty" gorm:"foreignKey:AssignedGuardID;references:ID"`
	Incident *Incident `json:"incident,omitempty" gorm:"foreignKey:AlertID;references:ID"`
}

type AlertType string
const (
	AlertTypeUnauthorizedAccess AlertType = "unauthorized_access"
	AlertTypeSuspiciousActivity AlertType = "suspicious_activity"
	AlertTypeEquipmentDamage    AlertType = "equipment_damage"
	AlertTypeSystemFailure      AlertType = "system_failure"
)

type AlertSeverity string
const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

type AlertStatus string
const (
	AlertStatusPending     AlertStatus = "pending"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusAssigned    AlertStatus = "assigned"
	AlertStatusResolved    AlertStatus = "resolved"
	AlertStatusClosed      AlertStatus = "closed"
)

type Incident struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AlertID        uuid.UUID      `json:"alert_id" gorm:"type:uuid;unique;not null"`
	Status         IncidentStatus `json:"status" gorm:"default:'open'"`
	Location       string         `json:"location" gorm:"not null"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// Relationships
	Alert          Alert            `json:"alert,omitempty" gorm:"foreignKey:AlertID;references:ID"`
	AssignedGuards []User           `json:"assigned_guards,omitempty" gorm:"many2many:incident_guards;joinForeignKey:IncidentID;JoinReferences:GuardID"`
	Updates        []IncidentUpdate `json:"updates,omitempty" gorm:"foreignKey:IncidentID;references:ID"`
}

type IncidentStatus string
const (
	IncidentStatusOpen        IncidentStatus = "open"
	IncidentStatusInProgress  IncidentStatus = "in_progress"
	IncidentStatusResolved    IncidentStatus = "resolved"
	IncidentStatusClosed      IncidentStatus = "closed"
)

type IncidentUpdate struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IncidentID uuid.UUID   `json:"incident_id" gorm:"type:uuid;not null"`
	GuardID    uuid.UUID   `json:"guard_id" gorm:"type:uuid;not null"`
	Type       UpdateType  `json:"type" gorm:"not null"`
	Message    string      `json:"message" gorm:"not null"`
    MediaURLs  pq.StringArray `json:"media_urls,omitempty" gorm:"type:text[]" swaggertype:"array,string"`
	Location   string      `json:"location,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`

	// Relationships
	Incident Incident `json:"incident,omitempty" gorm:"foreignKey:IncidentID;references:ID"`
	Guard    User     `json:"guard,omitempty" gorm:"foreignKey:GuardID;references:ID"`
}

type UpdateType string
const (
	UpdateTypeArrival       UpdateType = "arrival"
	UpdateTypeInvestigation UpdateType = "investigation"
	UpdateTypeResolution    UpdateType = "resolution"
)

// =======================
// Custom Join Tables
// =======================

type CameraGuard struct {
	CameraID uuid.UUID `json:"camera_id" gorm:"type:uuid;primaryKey"`
	GuardID  uuid.UUID `json:"guard_id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Camera Camera `json:"camera,omitempty" gorm:"foreignKey:CameraID;references:ID"`
	Guard  User   `json:"guard,omitempty" gorm:"foreignKey:GuardID;references:ID"`
}

type IncidentGuard struct {
	IncidentID uuid.UUID `json:"incident_id" gorm:"type:uuid;primaryKey"`
	GuardID    uuid.UUID `json:"guard_id" gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Incident Incident `json:"incident,omitempty" gorm:"foreignKey:IncidentID;references:ID"`
	Guard    User     `json:"guard,omitempty" gorm:"foreignKey:GuardID;references:ID"`
}



// BeforeCreate hook to set UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (p *Premise) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (c *Camera) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (i *Incident) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (iu *IncidentUpdate) BeforeCreate(tx *gorm.DB) error {
	if iu.ID == uuid.Nil {
		iu.ID = uuid.New()
	}
	return nil
} 