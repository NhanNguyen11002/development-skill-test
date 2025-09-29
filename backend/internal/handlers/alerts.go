package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/handlers/dto"
	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
)

// AlertHandler handles alert-related requests
type AlertHandler struct {
	service services.AlertsService
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler(service services.AlertsService) *AlertHandler {
	return &AlertHandler{
		service: service,
	}
}

// GetAlerts godoc
// @Summary Get alerts
// @Description Get all alerts with optional filters
// @Tags alerts
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param severity query string false "Filter by severity"
// @Param type query string false "Filter by type"
// @Param premise_id query string false "Filter by premise ID"
// @Success 200 {array} models.Alert
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts [get]
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	filters := services.AlertsFilter{
		Status:    c.Query("status"),
		Severity:  c.Query("severity"),
		Type:      c.Query("type"),
		PremiseID: c.Query("premise_id"),
	}

	role, _ := c.Get("role")
	userID := c.GetString("user_id")

	alerts, err := h.service.GetAlerts(c.Request.Context(), filters, role.(models.UserRole), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, alerts)
}

// GetAlert godoc
// @Summary Get alert by ID
// @Description Get alert details by ID
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} models.Alert
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts/{id} [get]
func (h *AlertHandler) GetAlert(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")
	userID := c.GetString("user_id")

	alert, err := h.service.GetAlert(c.Request.Context(), id, role.(models.UserRole), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Alert not found", err)
		return
	}
	response.Success(c, http.StatusOK, alert)
}

// AcknowledgeAlert godoc
// @Summary Acknowledge alert
// @Description Acknowledge the specified alert (SCS Operator only)
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} models.Alert
// @Failure 403 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts/{id}/acknowledge [post]
func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get("role")

	alert, err := h.service.AcknowledgeAlert(c.Request.Context(), id, role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, alert)
}

// AssignAlert godoc
// @Summary Assign alert to guard
// @Description Assign an alert to a security guard (SCS Operator only)
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param payload body dto.AssignAlertRequest true "Assign payload"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.ApiResponse
// @Failure 403 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts/{id}/assign [post]
func (h *AlertHandler) AssignAlert(c *gin.Context) {
	id := c.Param("id")
	var req dto.AssignAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	alert, incident, err := h.service.AssignAlert(c.Request.Context(), id, req.GuardID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"alert":    alert,
		"incident": incident,
	})
}

// CreateAlert godoc
// @Summary Create alert
// @Description Create a new alert (testing/demo)
// @Tags alerts
// @Accept json
// @Produce json
// @Param payload body dto.CreateAlertRequest true "Create alert payload"
// @Success 201 {object} models.Alert
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var alert models.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	created, err := h.service.CreateAlert(c.Request.Context(), alert)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create alert", err)
		return
	}
	response.Success(c, http.StatusCreated, created)
}

// UpdateAlert godoc
// @Summary Update alert status
// @Description Update an alert's status (SCS Operator only)
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param payload body dto.UpdateAlertStatusRequest true "Update alert status"
// @Success 200 {object} models.Alert
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/alerts/{id} [put]
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	alert, err := h.service.UpdateAlert(c.Request.Context(), id, models.AlertStatus(req.Status))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update alert", err)
		return
	}
	response.Success(c, http.StatusOK, alert)
}

