package controllers

import (
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	attendanceService *services.AttendanceService
}

func NewAttendanceController() *AttendanceController {
	return &AttendanceController{
		attendanceService: services.NewAttendanceService(),
	}
}

// ClockIn godoc
// @Summary Clock in
// @Description Record employee clock in
// @Tags attendance
// @Accept json
// @Produce json
// @Param attendance body models.AttendanceRequest true "Clock in data"
// @Success 201 {object} utils.Response{data=models.Attendance}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /attendance/clock-in [post]
func (c *AttendanceController) ClockIn(ctx *gin.Context) {
	var req models.AttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	attendance, err := c.attendanceService.ClockIn(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "Clock in successful", attendance.ToResponse())
}

// ClockOut godoc
// @Summary Clock out
// @Description Record employee clock out
// @Tags attendance
// @Accept json
// @Produce json
// @Param attendance body models.ClockOutRequest true "Clock out data"
// @Success 200 {object} utils.Response{data=models.Attendance}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /attendance/clock-out [put]
func (c *AttendanceController) ClockOut(ctx *gin.Context) {
	var req models.ClockOutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	attendance, err := c.attendanceService.ClockOut(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Clock out successful", attendance.ToResponse())
}

// GetAttendanceLogs godoc
// @Summary Get attendance logs
// @Description Get paginated attendance logs with filtering options
// @Tags attendance
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Param employee_id query string false "Filter by employee ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=[]models.AttendanceResponse}
// @Failure 500 {object} utils.Response
// @Router /attendance/logs [get]
func (c *AttendanceController) GetAttendanceLogs(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	employeeID := ctx.Query("employee_id")
	
	departmentID, _ := strconv.ParseUint(ctx.Query("department_id"), 10, 32)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	attendances, pagination, err := c.attendanceService.GetAttendanceLogs(
		startDate, endDate, uint(departmentID), employeeID, page, limit,
	)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Convert to response objects
	var attendanceResponses []models.AttendanceResponse
	for _, attendance := range attendances {
		attendanceResponses = append(attendanceResponses, attendance.ToResponse())
	}

	response := map[string]interface{}{
		"attendances": attendanceResponses,
		"pagination":  pagination,
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Attendance logs retrieved successfully", response)
}

// GetEmployeeAttendance godoc
// @Summary Get employee attendance
// @Description Get attendance records for a specific employee
// @Tags attendance
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response{data=[]models.AttendanceResponse}
// @Failure 500 {object} utils.Response
// @Router /attendance/employee/{employee_id} [get]
func (c *AttendanceController) GetEmployeeAttendance(ctx *gin.Context) {
	employeeID := ctx.Param("employee_id")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	attendances, err := c.attendanceService.GetEmployeeAttendance(employeeID, startDate, endDate)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	var attendanceResponses []models.AttendanceResponse
	for _, attendance := range attendances {
		attendanceResponses = append(attendanceResponses, attendance.ToResponse())
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee attendance retrieved successfully", attendanceResponses)
}

// GetAttendanceStats godoc
// @Summary Get attendance statistics
// @Description Get monthly attendance statistics for an employee
// @Tags attendance
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 500 {object} utils.Response
// @Router /attendance/stats/{employee_id} [get]
func (c *AttendanceController) GetAttendanceStats(ctx *gin.Context) {
	employeeID := ctx.Param("employee_id")
	month, _ := strconv.Atoi(ctx.Query("month"))
	year, _ := strconv.Atoi(ctx.Query("year"))

	stats, err := c.attendanceService.GetAttendanceStats(employeeID, month, year)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Attendance statistics retrieved successfully", stats)
}