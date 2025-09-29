package handlers

import (
	"net/http"

	"smart-city-surveillance/internal/config"
	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// AuthHandler handles authentication requests
type AuthHandler struct {
	config  *config.Config
	service services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(config *config.Config, service services.AuthService) *AuthHandler {
	return &AuthHandler{
		config:  config,
		service: service,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 401 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// Hide exact cause for security
		response.Error(c, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	resp := LoginResponse{
		Token: token,
		User:  user,
	}

	response.Success(c, http.StatusOK, resp)
}

// Logout godoc
// @Summary User logout
// @Description Logout current user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get currently authenticated user details
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Security BearerAuth
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	response.Success(c, http.StatusOK, user)
}
