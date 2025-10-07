// config/config.go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	// App
	AppEnv  string
	AppPort string
	GinMode string
	
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DatabaseURL string // For Railway
	
	// JWT
	JWTSecret string
	JWTExpiry string
	
	// CORS
	CORSAllowOrigin  string
	CORSAllowMethods string
	CORSAllowHeaders string
	
	// Email
	ResendAPIKey string
	FromEmail    string
	FrontendURL  string
}

func GetConfig() *Config {
	// For Railway, DATABASE_URL will be provided
	databaseURL := os.Getenv("DATABASE_URL")
	
	// Set default values
	appPort := os.Getenv("PORT") // Railway uses PORT
	if appPort == "" {
		appPort = "8080"
	}
	
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "production" // Default to production on Railway
	}
	
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release" // Default to release on Railway
	}

	return &Config{
		// App
		AppEnv:  appEnv,
		AppPort: appPort,
		GinMode: ginMode,
		
		// Database - Support both individual vars and DATABASE_URL
		DBHost:      getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:      getEnvWithDefault("DB_PORT", "3306"),
		DBUser:      getEnvWithDefault("DB_USER", "root"),
		DBPassword:  getEnvWithDefault("DB_PASSWORD", ""),
		DBName:      getEnvWithDefault("DB_NAME", "attendance_system"),
		DatabaseURL: databaseURL,
		
		// JWT
		JWTSecret: getEnvWithDefault("JWT_SECRET", "super-secret-jwt-key-here"),
		JWTExpiry: getEnvWithDefault("JWT_EXPIRY", "24h"),
		
		// CORS
		CORSAllowOrigin:  getEnvWithDefault("CORS_ALLOW_ORIGIN", "*"),
		CORSAllowMethods: getEnvWithDefault("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowHeaders: getEnvWithDefault("CORS_ALLOW_HEADERS", "*"),
		
		// Email
		ResendAPIKey: getEnvWithDefault("RESEND_API_KEY", ""),
		FromEmail:    getEnvWithDefault("FROM_EMAIL", "Attendance System <noreply@yourdomain.com>"),
		FrontendURL:  getEnvWithDefault("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}