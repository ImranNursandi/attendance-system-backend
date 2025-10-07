package main

import (
	"attendance-system/database"
	"attendance-system/middleware"
	"attendance-system/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Attendance System API
// @version 1.0
// @description API for Employee Attendance Management System

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	log.Println("🔍 Environment Debug:")
	log.Printf("   PORT: %s", os.Getenv("PORT"))
	log.Printf("   DATABASE_URL: %s", maskURL(os.Getenv("DATABASE_URL")))
	log.Printf("   JWT_SECRET: %s", maskSecret(os.Getenv("JWT_SECRET")))
	log.Printf("   FRONTEND_URL: %s", os.Getenv("FRONTEND_URL"))
	log.Printf("   GIN_MODE: %s", os.Getenv("GIN_MODE"))
	log.Printf("   RUN_MIGRATIONS: %s", os.Getenv("RUN_MIGRATIONS"))
	log.Printf("   APP_ENV: %s", os.Getenv("APP_ENV"))
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		} else {
			log.Println("✅ Loaded .env file for local development")
		}
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Run migrations on startup
	if os.Getenv("RUN_MIGRATIONS") == "true" || os.Getenv("APP_ENV") == "production" {
		if err := database.RunMigrations(database.GetDB()); err != nil {
			log.Fatal("❌ Failed to run migrations:", err)
		}
	}

	// Create Gin router
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS())
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	// log.Printf("📚 API Documentation available at http://localhost:%s/api/v1/docs", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}