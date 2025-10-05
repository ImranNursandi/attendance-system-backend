package routes

import (
	"attendance-system/controllers"
	"attendance-system/middleware"
	"attendance-system/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// Register custom validators
	middleware.RegisterCustomValidators()

	// Initialize services and controllers
	authService := services.NewAuthService()
	authController := controllers.NewAuthController()
	employeeController := controllers.NewEmployeeController()
	departmentController := controllers.NewDepartmentController()
	attendanceController := controllers.NewAttendanceController()
	reportController := controllers.NewReportController()

	// API v1 group
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Public routes (no authentication required)
		public := api.Group("/auth")
		{
			public.POST("/login", authController.Login)
			public.POST("/register", authController.Register)
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Auth routes
			auth := protected.Group("/auth")
			{
				auth.GET("/profile", authController.GetProfile)
				auth.PUT("/profile", authController.UpdateProfile)
				auth.PUT("/change-password", authController.ChangePassword)
				auth.POST("/refresh", authController.RefreshToken)
			}

			// Employee routes
			employees := protected.Group("/employees")
			employees.Use(middleware.RoleMiddleware([]string{"employee", "manager", "admin"}))
			{
				employees.GET("", employeeController.GetAllEmployees)
				employees.GET("/search", employeeController.SearchEmployees)
				employees.GET("/department/:department_id", employeeController.GetEmployeesByDepartment)
				employees.GET("/:id", employeeController.GetEmployeeByID)
				
				// Use a different path structure to avoid conflicts
				employees.GET("/employee/:employee_id/details", employeeController.GetEmployeeWithUser)
				
				// Admin/Manager only routes
				adminEmployees := employees.Group("")
				adminEmployees.Use(middleware.RoleMiddleware([]string{"manager", "admin"}))
				{
					adminEmployees.POST("", employeeController.CreateEmployee)
					adminEmployees.PUT("/:id", employeeController.UpdateEmployee)
					adminEmployees.DELETE("/:id", employeeController.DeleteEmployee)
				}
			}

			// Department routes
			departments := protected.Group("/departments")
			departments.Use(middleware.RoleMiddleware([]string{"manager", "admin"}))
			{
				departments.GET("", departmentController.GetAllDepartments)
				departments.GET("/active", departmentController.GetActiveDepartments)
				departments.GET("/:id", departmentController.GetDepartmentByID)
				departments.POST("", departmentController.CreateDepartment)
				departments.PUT("/:id", departmentController.UpdateDepartment)
				departments.DELETE("/:id", departmentController.DeleteDepartment)
			}

			// Attendance routes
			attendance := protected.Group("/attendance")
			attendance.Use(middleware.RoleMiddleware([]string{"employee", "manager", "admin"}))
			{
				attendance.POST("/clock-in", attendanceController.ClockIn)
				attendance.PUT("/clock-out", attendanceController.ClockOut)
				attendance.GET("/logs", attendanceController.GetAttendanceLogs)
				attendance.GET("/employee/:employee_id", attendanceController.GetEmployeeAttendance)
				attendance.GET("/stats/:employee_id", attendanceController.GetAttendanceStats)
			}

			// Report routes (Manager and Admin only)
			reports := protected.Group("/reports")
			reports.Use(middleware.RoleMiddleware([]string{"manager", "admin"}))
			{
				reports.GET("/attendance", reportController.GenerateAttendanceReport)
				reports.GET("/summary", reportController.GenerateSummaryReport)
				reports.GET("/department/:department_id", reportController.GenerateDepartmentReport)
				
				// Export routes
				reports.GET("/export/attendance", reportController.ExportAttendanceReport)
				reports.GET("/export/summary", reportController.ExportSummaryReport)
				reports.GET("/export/department/:department_id", reportController.ExportDepartmentReport)
			}

			// Dashboard routes
			protected.GET("/admin/dashboard", middleware.RoleMiddleware([]string{"admin"}), adminDashboard)
			protected.GET("/manager/dashboard", middleware.RoleMiddleware([]string{"manager", "admin"}), managerDashboard)
		}
	}

	// 404 handler
	router.NoRoute(notFoundHandler)
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"timestamp": time.Now().UTC(),
		"service":   "Attendance System API",
		"version":   "1.0.0",
	})
}

// adminDashboard handles admin dashboard data
func adminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin dashboard data",
		"data": gin.H{
			"total_employees":   150,
			"total_departments": 8,
			"today_attendance":  120,
			"pending_requests":  5,
		},
	})
}

// managerDashboard handles manager dashboard data
func managerDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Manager dashboard data",
		"data": gin.H{
			"department_employees": 25,
			"today_present":        22,
			"today_absent":         3,
			"pending_approvals":    2,
		},
	})
}

// notFoundHandler handles 404 errors
func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "Endpoint not found",
		"message": "The requested endpoint does not exist",
		"path":    c.Request.URL.Path,
	})
}