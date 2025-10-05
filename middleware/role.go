package middleware

import (
	"attendance-system/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if user has one of the required roles
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorJSON(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		role := userRole.(string)
		hasAccess := false

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			utils.ErrorJSON(c, http.StatusForbidden, 
				"Insufficient permissions. Required roles: "+strings.Join(allowedRoles, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware - only allows admin users
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"admin"})
}

// ManagerMiddleware - allows manager and admin users
func ManagerMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"manager", "admin"})
}

// EmployeeMiddleware - allows all authenticated users (employee, manager, admin)
func EmployeeMiddleware() gin.HandlerFunc {
	return RoleMiddleware([]string{"employee", "manager", "admin"})
}