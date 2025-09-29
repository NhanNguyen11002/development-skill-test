package dto
type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=active inactive maintenance"`
}
