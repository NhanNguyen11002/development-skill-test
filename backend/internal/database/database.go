package database

import (
	"fmt"
	"log"

	"smart-city-surveillance/internal/config"
	"smart-city-surveillance/internal/middleware"
	"smart-city-surveillance/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(config *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	log.Println("Running database migrations...")
	
	err := DB.AutoMigrate(
		&models.User{},
		&models.Premise{},
		&models.Camera{},
		&models.Alert{},
		&models.Incident{},
		&models.IncidentUpdate{},
		&models.CameraGuard{},
		&models.IncidentGuard{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedData populates the database with initial data
func SeedData() error {
	log.Println("Seeding database with initial data...")
	

	// Check if data already exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}
	pass,err := middleware.HashPassword("password")
	if(err != nil){
		return fmt.Errorf("failed to hash password")
	}
	// Create sample users
	users := []models.User{
		{
			Username:  "operator1",
			Email:     "operator1@st-engineering.com",
			Password:  pass, // password
			Role:      models.RoleSCSOperator,
			FirstName: "John",
			LastName:  "Operator",
			Phone:     "+65 9123 4567",
			IsActive:  true,
		},
		{
			Username:  "guard1",
			Email:     "guard1@st-engineering.com",
			Password:  pass, // password
			Role:      models.RoleSecurityGuard,
			FirstName: "Mike",
			LastName:  "Guard",
			Phone:     "+65 9123 4568",
			IsActive:  true,
		},
		{
			Username:  "guard2",
			Email:     "guard2@st-engineering.com",
			Password:  pass, // password
			Role:      models.RoleSecurityGuard,
			FirstName: "Sarah",
			LastName:  "Guard",
			Phone:     "+65 9123 4569",
			IsActive:  true,
		},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Username, err)
		}
	}

	// Create sample premises
	premises := []*models.Premise{
		{
			Name:        "ST Engineering HQ",
			Address:     "1 Ang Mo Kio Electronics Park Road, Singapore 567710",
			Type:        models.PremiseTypeOffice,
			FloorPlans:  "https://example.com/floorplans/hq.pdf",
			Description: "Main headquarters building",
			IsActive:    true,
		},
		{
			Name:        "Jurong Substation",
			Address:     "50 Jurong West Street 93, Singapore 648965",
			Type:        models.PremiseTypeSubstation,
			FloorPlans:  "https://example.com/floorplans/jurong.pdf",
			Description: "Power grid substation",
			IsActive:    true,
		},
		{
			Name:        "Woodlands Substation",
			Address:     "30 Woodlands Avenue 2, Singapore 738343",
			Type:        models.PremiseTypeSubstation,
			FloorPlans:  "https://example.com/floorplans/woodlands.pdf",
			Description: "Power grid substation",
			IsActive:    true,
		},
	}

	for index, premise := range premises {
		result :=  DB.Create(&premise)
		if result.Error != nil {
			return fmt.Errorf("failed to create premise %s: %w", premise.Name, result.Error)
		}
		premises[index].ID = premise.ID 
	}
	for i, p := range premises {
		fmt.Printf("Premise %d: name=%s, ID=%v\n", i, p.Name, p.ID)
	}

	// Create sample cameras
	cameras := []models.Camera{
		{
			Name:      "Main Entrance",
			Location:  "Front Gate",
			StreamURL: "rtsp://camera1.example.com/stream1",
			Status:    models.CameraStatusActive,
			PremiseID: premises[0].ID, // ST Engineering HQ
		},
		{
			Name:      "Parking Lot",
			Location:  "Underground Parking",
			StreamURL: "rtsp://camera2.example.com/stream2",
			Status:    models.CameraStatusActive,
			PremiseID: premises[0].ID, // ST Engineering HQ
		},
		{
			Name:      "Equipment Room",
			Location:  "Basement Level 1",
			StreamURL: "rtsp://camera3.example.com/stream3",
			Status:    models.CameraStatusActive,
			PremiseID: premises[1].ID, // Jurong Substation
		},
		{
			Name:      "Perimeter Fence",
			Location:  "North Side",
			StreamURL: "rtsp://camera4.example.com/stream4",
			Status:    models.CameraStatusActive,
			PremiseID: premises[1].ID, // Jurong Substation
		},
		{
			Name:      "Control Room",
			Location:  "Main Building",
			StreamURL: "rtsp://camera5.example.com/stream5",
			Status:    models.CameraStatusActive,
			PremiseID: premises[2].ID, // Woodlands Substation
		},
		{
			Name:      "Generator Area",
			Location:  "Back Yard",
			StreamURL: "rtsp://camera6.example.com/stream6",
			Status:    models.CameraStatusActive,
			PremiseID: premises[2].ID, // Woodlands Substation
		},
	}

	for _, camera := range cameras {
		if err := DB.Create(&camera).Error; err != nil {
			return fmt.Errorf("failed to create camera %s: %w", camera.Name, err)
		}
	}

	// Assign guards to cameras
	var guards []models.User
	DB.Where("role = ?", models.RoleSecurityGuard).Find(&guards)

	var allCameras []models.Camera
	DB.Find(&allCameras)

	// Assign guard1 to first 3 cameras, guard2 to last 3 cameras
	for i, camera := range allCameras {
		guardIndex := i / 3
		if guardIndex < len(guards) {
			cameraGuard := models.CameraGuard{
				CameraID: camera.ID,
				GuardID:  guards[guardIndex].ID,
			}
			if err := DB.Create(&cameraGuard).Error; err != nil {
				return fmt.Errorf("failed to assign guard to camera: %w", err)
			}
		}
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
} 