package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
)

// IncidentHandler handles incident-related requests
type IncidentHandler struct {
	service services.IncidentsService
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler(service services.IncidentsService) *IncidentHandler {
	return &IncidentHandler{
		service: service,
	}
}

// GetIncidents godoc
// @Summary Get incidents
// @Description Get incidents, filtered by status; guards see only their incidents
// @Tags incidents
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Incident
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents [get]
func (h *IncidentHandler) GetIncidents(c *gin.Context) {
	userRole, _ := c.Get("role")
	userID := c.GetString("user_id")
	status := c.Query("status")

	incidents, err := h.service.GetIncidents(c.Request.Context(), userRole.(models.UserRole), userID, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch incidents", err)
		return
	}
	response.Success(c, http.StatusOK, incidents)
}

// GetIncident godoc
// @Summary Get incident by ID
// @Description Get incident details by ID; guards must be assigned
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path string true "Incident ID"
// @Success 200 {object} models.Incident
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents/{id} [get]
func (h *IncidentHandler) GetIncident(c *gin.Context) {
	id := c.Param("id")
	userRole, _ := c.Get("role")
	userID := c.GetString("user_id")

	incident, err := h.service.GetIncident(c.Request.Context(), id, userRole.(models.UserRole), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Incident not found", err)
		return
	}
	response.Success(c, http.StatusOK, incident)
}

// UpdateIncident godoc
// @Summary Update incident status
// @Description Update an incident's status; guards must be assigned
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path string true "Incident ID"
// @Param payload body dto.UpdateAlertStatusRequest true "Update payload"
// @Success 200 {object} models.Incident
// @Failure 400 {object} response.ApiResponse
// @Failure 403 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents/{id} [put]
func (h *IncidentHandler) UpdateIncident(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userRole, _ := c.Get("role")
	userID := c.GetString("user_id")

	incident, err := h.service.UpdateIncident(c.Request.Context(), id, models.IncidentStatus(req.Status), userRole.(models.UserRole), userID)
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Access denied", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update incident", err)
		return
	}
	response.Success(c, http.StatusOK, incident)
}

// AddIncidentUpdate godoc
// @Summary Add incident update
// @Description Add an update to an incident; guards must be assigned
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path string true "Incident ID"
// @Param payload body models.IncidentUpdate true "Update payload"
// @Success 201 {object} models.IncidentUpdate
// @Failure 400 {object} response.ApiResponse
// @Failure 403 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents/{id}/updates [post]
func (h *IncidentHandler) AddIncidentUpdate(c *gin.Context) {
	id := c.Param("id")
	var update models.IncidentUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userRole, _ := c.Get("role")
	userID := c.GetString("user_id")

	saved, err := h.service.AddIncidentUpdate(c.Request.Context(), id, update, userRole.(models.UserRole), userID)
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Access denied", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create update", err)
		return
	}
	response.Success(c, http.StatusCreated, saved)
}

// GetIncidentByAlertID godoc
// @Summary Get incident by alert ID
// @Description Get an incident by its associated alert ID (SCS Operator)
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} models.Incident
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents/by-alert/{id} [get]
func (h *IncidentHandler) GetIncidentByAlertID(c *gin.Context) {
	id := c.Param("id")
	incident, err := h.service.GetIncidentByAlertID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Incident not found", err)
		return
	}
	response.Success(c, http.StatusOK, incident)
}

// GetAssignedIncidents godoc
// @Summary Get my assigned incidents
// @Description Get incidents assigned to the authenticated guard
// @Tags incidents
// @Accept json
// @Produce json
// @Success 200 {array} models.Incident
// @Failure 401 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/incidents/assigned/me [get]
func (h *IncidentHandler) GetAssignedIncidents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	incidents, err := h.service.GetIncidents(c.Request.Context(), models.RoleSecurityGuard, userID.(string), "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch incidents", err)
		return
	}
	response.Success(c, http.StatusOK, incidents)
} 