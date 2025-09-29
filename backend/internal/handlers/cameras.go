package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
)

// CameraHandler handles camera-related endpoints
type CameraHandler struct {
	service services.CameraService
}

func NewCameraHandler(service services.CameraService) *CameraHandler {
	return &CameraHandler{service: service}
}

// GetCameras godoc
// @Summary Get all cameras
// @Description Get a list of all cameras (SCS Operator only)
// @Tags cameras
// @Accept json
// @Produce json
// @Success 200 {array} models.Camera
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/cameras [get]
func (h *CameraHandler) GetCameras(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	cameras, err := h.service.GetAll(c.Request.Context(), role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, cameras)
}

// GetCamera godoc
// @Summary Get camera by ID
// @Description Get camera details by ID (SCS Operator or assigned Security Guard)
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path string true "Camera ID"
// @Success 200 {object} models.Camera
// @Failure 403 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/cameras/{id} [get]
func (h *CameraHandler) GetCamera(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}
	 userID := c.GetString("user_id")

	id := c.Param("id")
	camera, err := h.service.GetByID(c.Request.Context(), id, userID, role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError,"Internal Server", err)
		return
	}
	if camera == nil {
		response.Error(c, http.StatusNotFound, "Camera not found", nil)
		return
	}
	response.Success(c, http.StatusOK, camera)
}

// UpdateCameraStatus godoc
// @Summary Update camera status
// @Description Update the status of a camera (SCS Operator only)
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path string true "Camera ID"
// @Param status body dto.UpdateStatusRequest true "New status"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/cameras/{id}/status [put]
func (h *CameraHandler) UpdateCameraStatus(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive maintenance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Internal Server",err)
		return
	}

	id := c.Param("id")
	err := h.service.UpdateStatus(c.Request.Context(), id, models.CameraStatus(req.Status), role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server",err)
		return
	}
	response.Success(c, http.StatusOK, nil )
}

// GetCamerasByPremise godoc
// @Summary Get cameras by premise ID
// @Description Get cameras under a premise (SCS Operator only)
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path string true "Premise ID"
// @Success 200 {array} models.Camera
// @Failure 403 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/cameras/premise/{id} [get]
func (h *CameraHandler) GetCamerasByPremise(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role not found", nil)
		return
	}

	premiseID := c.Param("id")
	cameras, err := h.service.GetByPremiseID(c.Request.Context(), premiseID, role.(models.UserRole))
	if err != nil {
		if err.Error() == "permission denied" {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, cameras)
}

// GetAssignedCameras godoc
// @Summary Get cameras assigned to current guard
// @Description Get cameras assigned to the authenticated security guard
// @Tags cameras
// @Accept json
// @Produce json
// @Success 200 {array} models.Camera
// @Failure 401 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/cameras/assigned [get]
func (h *CameraHandler) GetAssignedCameras(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	cameras, err := h.service.GetAssignedByGuardID(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server", err)
		return
	}
	response.Success(c, http.StatusOK, cameras)
}
