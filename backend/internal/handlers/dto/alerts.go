package dto

type AssignAlertRequest struct {
	GuardID []string `json:"guard_id" binding:"required"`
}
type UpdateAlertStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending acknowledged assigned resolved closed"`

}

type CreateAlertRequest struct {
    Type        string `json:"type" binding:"required,oneof=unauthorized_access suspicious_activity equipment_damage system_failure"`
    Severity    string `json:"severity" binding:"required,oneof=low medium high critical"`
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
    Location    string `json:"location" binding:"required"`
    PremiseID   string `json:"premise_id" binding:"required,uuid"`
    CameraID    string `json:"camera_id,omitempty" binding:"omitempty,uuid"`
}