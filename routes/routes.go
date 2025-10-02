package routes

import (
	"attendance-system/config"
	"attendance-system/controllers"
	"attendance-system/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Initialize repositories
	employeeRepo := repositories.NewEmployeeRepository(config.DB)
	departmentRepo := repositories.NewDepartmentRepository(config.DB)
	attendanceRepo := repositories.NewAttendanceRepository(config.DB)

	// Initialize controllers
	employeeController := controllers.NewEmployeeController(employeeRepo)
	departmentController := controllers.NewDepartmentController(departmentRepo)
	attendanceController := controllers.NewAttendanceController(attendanceRepo, employeeRepo, departmentRepo)

	// Employee routes
	employeeRoutes := router.Group("/employees")
	{
		employeeRoutes.POST("", employeeController.CreateEmployee)
		employeeRoutes.GET("", employeeController.GetEmployees)
		employeeRoutes.GET("/:id", employeeController.GetEmployee)
		employeeRoutes.PUT("/:id", employeeController.UpdateEmployee)
		employeeRoutes.DELETE("/:id", employeeController.DeleteEmployee)
	}

	// Department routes
	deptRoutes := router.Group("/departments")
	{
		deptRoutes.POST("", departmentController.CreateDepartment)
		deptRoutes.GET("", departmentController.GetDepartments)
		deptRoutes.GET("/:id", departmentController.GetDepartment)
		deptRoutes.PUT("/:id", departmentController.UpdateDepartment)
		deptRoutes.DELETE("/:id", departmentController.DeleteDepartment)
	}

	// Attendance routes
	attendanceRoutes := router.Group("/attendance")
	{
		attendanceRoutes.POST("/clock-in", attendanceController.ClockIn)
		attendanceRoutes.PUT("/clock-out", attendanceController.ClockOut)
		attendanceRoutes.GET("/logs", attendanceController.GetAttendanceLogs)
	}

	return router
}