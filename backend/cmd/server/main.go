package main

// @title My API
// @version 1.0
// @description This is my API.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"log"
	"net/http"

	_ "smart-city-surveillance/docs"
	"smart-city-surveillance/internal/config"
	"smart-city-surveillance/internal/database"
	"smart-city-surveillance/internal/handlers"
	"smart-city-surveillance/internal/middleware"
	"smart-city-surveillance/internal/models"
	"smart-city-surveillance/internal/services"
	"smart-city-surveillance/pkg/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Smart City Surveillance API
// @version 1.0
// @description API server cho Smart City Surveillance project

// @host localhost:8081
// @BasePath /

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed data
	if err := database.SeedData(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize handlers
	//premises
	premisesService := services.NewPremisesService(database.GetDB())
	premiseHandler := handlers.NewPremiseHandler(premisesService)

	camerasService := services.NewCameraService(database.GetDB())
	cameraHandler := handlers.NewCameraHandler(camerasService)

		// Users
	userService := services.NewUserService(database.GetDB())
	userHandler := handlers.NewUserHandler(userService)

	// Auth
	authService := services.NewAuthService(database.GetDB(), cfg)
	authHandler := handlers.NewAuthHandler(cfg, authService)

	// Alerts
	alertsService := services.NewAlertsService(database.GetDB(), wsHub)
	alertHandler := handlers.NewAlertHandler(alertsService)

	// Incidents
	incidentsService := services.NewIncidentsService(database.GetDB(), wsHub)
	incidentHandler := handlers.NewIncidentHandler(incidentsService)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5175","http://localhost:5174","http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "smart-city-surveillance",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetCurrentUser)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
							// Premises routes
				premises := protected.Group("/premises")
				{
					premises.GET("", middleware.RoleMiddleware(models.RoleSCSOperator), premiseHandler.GetPremises)
					premises.GET("/:id", middleware.RoleMiddleware(models.RoleSCSOperator), premiseHandler.GetPremise)
					premises.GET("/:id/cameras", premiseHandler.GetPremiseCameras)
				}

						// Cameras routes
				cameras := protected.Group("/cameras")
				{
					cameras.GET("", cameraHandler.GetCameras)
					cameras.GET("/premise/:id", middleware.RoleMiddleware(models.RoleSCSOperator), cameraHandler.GetCamerasByPremise)
					cameras.GET("/assigned", middleware.RoleMiddleware(models.RoleSecurityGuard), cameraHandler.GetAssignedCameras)
					cameras.GET("/:id", cameraHandler.GetCamera)
					cameras.PUT("/:id/status", middleware.RoleMiddleware(models.RoleSCSOperator), cameraHandler.UpdateCameraStatus)
				}

							// Alerts routes
				alerts := protected.Group("/alerts")
				{
					alerts.GET("", middleware.RoleMiddleware(models.RoleSCSOperator), alertHandler.GetAlerts)
					alerts.GET("/:id", alertHandler.GetAlert)
					alerts.POST("/:id/acknowledge", middleware.RoleMiddleware(models.RoleSCSOperator), alertHandler.AcknowledgeAlert)
					alerts.POST("/:id/assign", middleware.RoleMiddleware(models.RoleSCSOperator), alertHandler.AssignAlert)
					alerts.POST("", middleware.RoleMiddleware(models.RoleSCSOperator), alertHandler.CreateAlert)
					alerts.PUT("/:id", middleware.RoleMiddleware(models.RoleSCSOperator), alertHandler.UpdateAlert)
				}

							// Incidents routes
				incidents := protected.Group("/incidents")
				{
					incidents.GET("", middleware.RoleMiddleware(models.RoleSCSOperator), incidentHandler.GetIncidents)
					incidents.GET("/:id", middleware.RoleMiddleware(models.RoleSCSOperator), incidentHandler.GetIncident)
					incidents.GET("/by-alert/:id", middleware.RoleMiddleware(models.RoleSCSOperator), incidentHandler.GetIncidentByAlertID)
					incidents.GET("/assigned/me", middleware.RoleMiddleware(models.RoleSecurityGuard), incidentHandler.GetAssignedIncidents)
					incidents.PUT("/:id", incidentHandler.UpdateIncident)
					incidents.POST("/:id/updates", incidentHandler.AddIncidentUpdate)
				}
				

						// Users routes
				users := protected.Group("/users")
				{
					users.GET("", userHandler.GetUsers)
					users.GET("/assigned/camera/:id", middleware.RoleMiddleware(models.RoleSCSOperator), userHandler.GetUsersByAssignedCamera)
					users.GET("/assigned/incident/:id", middleware.RoleMiddleware(models.RoleSCSOperator), userHandler.GetUsersByAssignedIncident)
				}
		}

		// WebSocket endpoint
		router.GET("/ws", func(c *gin.Context) {
			// Extract user info from query params (in real app, use JWT)
			userID := c.Query("user_id")
			role := c.Query("role")
			
			if userID == "" || role == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and role required"})
				return
			}

			websocket.ServeWebSocket(wsHub, userID, role)(c.Writer, c.Request)
		})
	}

	// Start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 