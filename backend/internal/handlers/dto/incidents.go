package dto

type UpdateIncidentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending acknowledged assigned resolved closed"`

}