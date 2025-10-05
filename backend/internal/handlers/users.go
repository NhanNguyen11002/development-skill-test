package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related endpoints

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users (SCS Operator only)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	users, err := h.service.GetAll(c.Request.Context(), role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, users)
}

// GetUsersByAssignedCamera godoc
// @Summary Get users assigned to a camera
// @Description Get users assigned to a specific camera (Operator only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Camera ID"
// @Success 200 {array} models.User
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/users/assigned/camera/{id} [get]
func (h *UserHandler) GetUsersByAssignedCamera(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	cameraID := c.Param("id")
	users, err := h.service.GetByAssignedCameraID(c.Request.Context(), cameraID, role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, users)
}

// GetUsersByAssignedIncident godoc
// @Summary Get users assigned to an incident
// @Description Get users assigned to a specific incident (Operator only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Incident ID"
// @Success 200 {array} models.User
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/users/assigned/incident/{id} [get]
func (h *UserHandler) GetUsersByAssignedIncident(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	incidentID := c.Param("id")
	users, err := h.service.GetByAssignedIncidentID(c.Request.Context(), incidentID, role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, users)
} 