package handlers

import (
	"attendance-system/config"
	"attendance-system/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Status     string                 `json:"status"`
	Timestamp  time.Time              `json:"timestamp"`
	Service    string                 `json:"service"`
	Version    string                 `json:"version"`
	Database   DatabaseStatus         `json:"database"`
	System     SystemInfo             `json:"system"`
	Uptime     string                 `json:"uptime"`
	Components map[string]interface{} `json:"components,omitempty"`
}

type DatabaseStatus struct {
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	Latency   string `json:"latency,omitempty"`
}

type SystemInfo struct {
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
	Hostname  string `json:"hostname,omitempty"`
}

var startTime = time.Now()

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := checkDatabaseStatus()

	// Prepare health response
	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "Attendance System API",
		Version:   "1.0.0",
		Database:  dbStatus,
		System: SystemInfo{
			GoVersion: "1.21",
			Platform:  "Golang",
		},
		Uptime: formatUptime(time.Since(startTime)),
		Components: map[string]interface{}{
			"api":        "healthy",
			"database":   dbStatus.Status,
			"middleware": "healthy",
		},
	}

	// If database is not connected, mark overall status as degraded
	if !dbStatus.Connected {
		health.Status = "degraded"
		health.Components["database"] = "unhealthy"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service health status",
		"data":    health,
	})
}

// DetailedHealthCheck provides detailed health information
func DetailedHealthCheck(c *gin.Context) {
	dbStatus := checkDatabaseStatus()

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "Attendance System API",
		Version:   "1.0.0",
		Database:  dbStatus,
		System: SystemInfo{
			GoVersion: "1.21",
			Platform:  "Golang",
		},
		Uptime: formatUptime(time.Since(startTime)),
		Components: map[string]interface{}{
			"api": map[string]interface{}{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
			},
			"database": map[string]interface{}{
				"status":    dbStatus.Status,
				"connected": dbStatus.Connected,
				"latency":   dbStatus.Latency,
			},
			"memory": map[string]interface{}{
				"status": "healthy",
				"usage":  "normal",
			},
		},
	}

	if !dbStatus.Connected {
		health.Status = "degraded"
		health.Components["database"].(map[string]interface{})["status"] = "unhealthy"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detailed service health status",
		"data":    health,
	})
}

// checkDatabaseStatus checks the database connection status
func checkDatabaseStatus() DatabaseStatus {
	start := time.Now()
	
	var result int
	err := config.DB.Raw("SELECT 1").Scan(&result).Error
	
	latency := time.Since(start).String()
	
	if err != nil {
		return DatabaseStatus{
			Status:    "unhealthy",
			Connected: false,
			Latency:   latency,
		}
	}

	return DatabaseStatus{
		Status:    "healthy",
		Connected: true,
		Latency:   latency,
	}
}

// formatUptime formats the uptime duration to a human-readable string
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return utils.FormatBytes(int64(days)) + " days"
	}
	if hours > 0 {
		return utils.FormatBytes(int64(hours)) + " hours"
	}
	if minutes > 0 {
		return utils.FormatBytes(int64(minutes)) + " minutes"
	}
	return utils.FormatBytes(int64(seconds)) + " seconds"
}

// ReadinessCheck checks if the service is ready to handle requests
func ReadinessCheck(c *gin.Context) {
	dbStatus := checkDatabaseStatus()

	if !dbStatus.Connected {
		utils.ErrorJSON(c, http.StatusServiceUnavailable, "Service not ready: Database unavailable")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
		},
	})
}

// LivenessCheck checks if the service is alive
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service is alive",
		"data": map[string]interface{}{
			"status":    "alive",
			"timestamp": time.Now().UTC(),
		},
	})
}