// config/config.go
package config

import (
	"os"
)

type Config struct {
	AppEnv       string
	AppPort      string
	GinMode      string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiry    string
	CORSOrigin   string
	CORSMethods  string
	CORSHeaders  string
	ResendAPIKey string
	FromEmail    string
	FrontendURL  string
	
	// Add individual DB config fields for local development
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var appConfig *Config

func GetConfig() *Config {
	if appConfig == nil {
		appConfig = &Config{
			AppEnv:       getEnv("APP_ENV", "development"),
			AppPort:      getEnv("PORT", "8080"),
			GinMode:      getEnv("GIN_MODE", "debug"),
			DatabaseURL:  Getenv("DATABASE_URL"),
			JWTSecret:    getEnv("JWT_SECRET", "super-secret-jwt-key-here"),
			JWTExpiry:    getEnv("JWT_EXPIRY", "24h"),
			CORSOrigin:   getEnv("CORS_ALLOW_ORIGIN", "*"),
			CORSMethods:  getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORSHeaders:  getEnv("CORS_ALLOW_HEADERS", "*"),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
			FromEmail:    getEnv("FROM_EMAIL", "onboarding@resend.dev"),
			FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
			
			// Individual DB config
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "3306"),
			DBUser:     getEnv("DB_USER", "root"),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "attendance_system"),
		}
	}
	return appConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}