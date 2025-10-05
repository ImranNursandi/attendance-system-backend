package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateID generates a random ID
func GenerateID(prefix string, length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return prefix + hex.EncodeToString(bytes)
}

// GenerateEmployeeID generates a unique employee ID
func GenerateEmployeeID() string {
	timestamp := time.Now().Unix()
	random, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return fmt.Sprintf("EMP%d%03d", timestamp, random)
}

// GenerateAttendanceID generates a unique attendance ID
func GenerateAttendanceID(employeeID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ATT%s%d", employeeID, timestamp)
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	// Simple phone validation - adjust based on your requirements
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove script tags and other potentially dangerous content
	input = strings.ReplaceAll(input, "<script>", "")
	input = strings.ReplaceAll(input, "</script>", "")
	input = strings.ReplaceAll(input, "javascript:", "")
	
	// Trim whitespace
	return strings.TrimSpace(input)
}

// GetPaginationParams extracts pagination parameters from request
func GetPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate and adjust values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

// GetDateRangeParams extracts date range parameters from request
func GetDateRangeParams(c *gin.Context) (string, string) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// If no dates provided, default to current month
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	return startDate, endDate
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}

// IsAdmin checks if user is admin
func IsAdmin(c *gin.Context) bool {
	return GetUserRoleFromContext(c) == "admin"
}

// GetClientIP gets the client IP address
func GetClientIP(c *gin.Context) string {
	// Check for forwarded IP first
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Fall back to remote address
	return c.ClientIP()
}

// FormatBytes formats bytes to human readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// CalculatePercentage calculates percentage
func CalculatePercentage(part, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

// GetCurrentTime returns current time in UTC
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// ParseTime parses time string with multiple formats
func ParseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// SendFileResponse sends file download response
func SendFileResponse(c *gin.Context, filename string, data []byte, contentType string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	
	c.Data(http.StatusOK, contentType, data)
}