package middleware

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware configuration
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()

	// Get allowed origins from environment or use default
	allowOrigins := os.Getenv("CORS_ALLOW_ORIGIN")
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	config.AllowOrigins = strings.Split(allowOrigins, ",")

	// Get allowed methods from environment or use default
	allowMethods := os.Getenv("CORS_ALLOW_METHODS")
	if allowMethods == "" {
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	} else {
		config.AllowMethods = strings.Split(allowMethods, ",")
	}

	// Get allowed headers from environment or use default
	allowHeaders := os.Getenv("CORS_ALLOW_HEADERS")
	if allowHeaders == "" {
		config.AllowHeaders = []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
		}
	} else {
		config.AllowHeaders = strings.Split(allowHeaders, ",")
	}

	// Get allow credentials from environment
	allowCredentials, _ := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	config.AllowCredentials = allowCredentials

	// Get max age from environment
	if maxAgeStr := os.Getenv("CORS_MAX_AGE"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
			config.MaxAge = time.Duration(maxAge) * time.Second
		}
	}

	return cors.New(config)
}