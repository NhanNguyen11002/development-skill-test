package dto

type UpdateIncidentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending acknowledged assigned resolved closed"`

}
type AddIncidentUpdateRequest struct {
	Type       string   `json:"type" binding:"required,oneof=arrival investigation resolution"`
	Message    string   `json:"message" binding:"required"`
	MediaURLs  []string `json:"media_urls,omitempty"`
	Location   string   `json:"location,omitempty"`
}