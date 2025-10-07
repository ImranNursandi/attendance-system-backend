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
	
	corsConfig := cors.Config{
		AllowMethods:     strings.Split(cfg.CORSMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORSHeaders, ","),
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	// Handle multiple origins
	if cfg.CORSOrigin == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = strings.Split(cfg.CORSOrigin, ",")
		corsConfig.AllowAllOrigins = false
	}

	return cors.New(corsConfig)
}