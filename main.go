package main

import (
	"attendance-system/database"
	"attendance-system/middleware"
	"attendance-system/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Debug: List ALL environment variables (be careful with secrets)
	log.Println("🔍 Checking ALL environment variables:")
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		// Mask sensitive values
		if strings.Contains(strings.ToLower(key), "secret") || 
		   strings.Contains(strings.ToLower(key), "password") ||
		   strings.Contains(strings.ToLower(key), "key") ||
		   strings.Contains(strings.ToLower(key), "token") {
			log.Printf("   %s: ***MASKED***", key)
		} else {
			log.Printf("   %s: %s", key, pair[1])
		}
	}

	// Try to load .env file (for local development only)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found - relying on system environment variables")
	} else {
		log.Println("✅ .env file loaded successfully")
	}

	// Now check our specific required variables
	log.Println("🔍 Required Environment Variables:")
	requiredVars := []string{"PORT", "DATABASE_URL", "JWT_SECRET", "FRONTEND_URL"}
	for _, v := range requiredVars {
		value := os.Getenv(v)
		if value == "" {
			log.Printf("   ❌ %s: NOT SET", v)
		} else {
			if v == "DATABASE_URL" {
				log.Printf("   ✅ %s: %s", v, maskURL(value))
			} else if v == "JWT_SECRET" {
				log.Printf("   ✅ %s: %s", v, maskSecret(value))
			} else {
				log.Printf("   ✅ %s: %s", v, value)
			}
		}
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("🚀 Running in RELEASE mode")
	} else {
		log.Println("🔧 Running in DEBUG mode")
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Run migrations on startup
	if os.Getenv("RUN_MIGRATIONS") == "true" {
		log.Println("🔄 Running database migrations...")
		if err := database.RunMigrations(database.GetDB()); err != nil {
			log.Fatal("❌ Failed to run migrations:", err)
		}
		log.Println("✅ Database migrations completed")
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

	log.Printf("🚀 Server starting on port %s", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func maskURL(url string) string {
	if url == "" {
		return "(not set)"
	}
	parts := strings.Split(url, "@")
	if len(parts) > 1 {
		return parts[0] + "@***"
	}
	return "***"
}

func maskSecret(secret string) string {
	if secret == "" {
		return "(not set)"
	}
	if len(secret) > 8 {
		return secret[:4] + "***" + secret[len(secret)-4:]
	}
	return "***"
}