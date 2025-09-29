package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PremisesHandler struct {
	service services.PremisesService
}

func NewPremiseHandler(service services.PremisesService) *PremisesHandler {
	return &PremisesHandler{service: service}
}

// GetPremises godoc
// @Summary Get all premises
// @Description Retrieve all premises in the system
// @Tags premises
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/premises [get]
func (h *PremisesHandler) GetPremises(c *gin.Context) {
	ctx := c.Request.Context()
	premises, err := h.service.GetPremises(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch premises", err)
		return
	}
	response.Success(c, http.StatusOK, premises)
}

// GetPremise godoc
// @Summary Get a premise by ID
// @Description Retrieve premise details by ID
// @Tags premises
// @Produce json
// @Param id path string true "Premise ID"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/premises/{id} [get]
func (h *PremisesHandler) GetPremise(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	idUUID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Premise not found", err)
			return
	}

	premise, err := h.service.GetPremise(ctx, idUUID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Premise not found", err)
		return
	}
	response.Success(c, http.StatusOK, premise)
}

// GetPremiseCameras godoc
// @Summary Get cameras in a premise
// @Description Retrieve all cameras located in a specific premise
// @Tags premises
// @Produce json
// @Param id path string true "Premise ID"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/premises/{id}/cameras [get]
func (h *PremisesHandler) GetPremiseCameras(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	idUUID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to fetch cameras", err)
			return
	}

	cameras, err := h.service.GetPremiseCameras(ctx, idUUID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to fetch cameras", err)
		return
	}
	response.Success(c, http.StatusOK, cameras)
}
