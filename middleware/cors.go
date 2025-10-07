// middleware/cors.go
package middleware

import (
	"attendance-system/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	cfg := config.GetConfig()
	
	config := cors.Config{
		AllowMethods:     strings.Split(cfg.CORSAllowMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORSAllowHeaders, ","),
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	// Handle multiple origins
	if cfg.CORSAllowOrigin == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = strings.Split(cfg.CORSAllowOrigin, ",")
		config.AllowAllOrigins = false
	}

	return cors.New(config)
}