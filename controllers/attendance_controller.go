package controllers

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceController struct {
	attendanceRepo *repositories.AttendanceRepository
	employeeRepo   *repositories.EmployeeRepository
	deptRepo       *repositories.DepartmentRepository
}

func NewAttendanceController(attendanceRepo *repositories.AttendanceRepository, employeeRepo *repositories.EmployeeRepository, deptRepo *repositories.DepartmentRepository) *AttendanceController {
	return &AttendanceController{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		deptRepo:       deptRepo,
	}
}

type ClockInRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
}

type ClockOutRequest struct {
	EmployeeID  string `json:"employee_id" binding:"required"`
	AttendanceID string `json:"attendance_id" binding:"required"`
}

func (c *AttendanceController) ClockIn(ctx *gin.Context) {
	var req ClockInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// Check if employee exists
	employee, err := c.employeeRepo.FindByID(req.EmployeeID)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Employee not found")
		return
	}

	// Check if already clocked in today
	existingAttendance, _ := c.attendanceRepo.GetAttendanceByEmployeeID(req.EmployeeID, time.Now())
	if existingAttendance.ID != "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Already clocked in today")
		return
	}

	now := time.Now()
	attendanceID := uuid.New().String()

	// Create attendance record
	attendance := models.Attendance{
		ID:         attendanceID,
		EmployeeID: req.EmployeeID,
		ClockIn:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Check if clock in is on time
	department, _ := c.deptRepo.FindByID(employee.DepartmentID)
	maxClockIn, _ := time.Parse("15:04", department.MaxClockIn)
	clockInTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	isOnTime := clockInTime.Before(maxClockIn) || clockInTime.Equal(maxClockIn)

	// Create attendance history
	history := models.AttendanceHistory{
		ID:             uuid.New().String(),
		EmployeeID:     req.EmployeeID,
		AttendanceID:   attendanceID,
		DateAttendance: now,
		AttendanceType: 1, // Clock In
		Description:    "Clock In - " + map[bool]string{true: "On Time", false: "Late"}[isOnTime],
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := c.attendanceRepo.ClockIn(&attendance, &history); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to clock in")
		return
	}

	response := gin.H{
		"attendance": attendance,
		"is_on_time": isOnTime,
		"message":    "Clock in successful",
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "Clock in recorded successfully", response)
}

func (c *AttendanceController) ClockOut(ctx *gin.Context) {
	var req ClockOutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// Check if employee exists
	employee, err := c.employeeRepo.FindByID(req.EmployeeID)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Employee not found")
		return
	}

	now := time.Now()

	// Check if clock out is on time
	department, _ := c.deptRepo.FindByID(employee.DepartmentID)
	maxClockOut, _ := time.Parse("15:04", department.MaxClockOut)
	clockOutTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	isOnTime := clockOutTime.After(maxClockOut) || clockOutTime.Equal(maxClockOut)

	// Create attendance history
	history := models.AttendanceHistory{
		ID:             uuid.New().String(),
		EmployeeID:     req.EmployeeID,
		AttendanceID:   req.AttendanceID,
		DateAttendance: now,
		AttendanceType: 2, // Clock Out
		Description:    "Clock Out - " + map[bool]string{true: "On Time", false: "Early"}[isOnTime],
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := c.attendanceRepo.ClockOut(req.AttendanceID, now, &history); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to clock out")
		return
	}

	response := gin.H{
		"is_on_time": isOnTime,
		"message":    "Clock out successful",
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Clock out recorded successfully", response)
}

func (c *AttendanceController) GetAttendanceLogs(ctx *gin.Context) {
	date := ctx.Query("date")
	deptID := ctx.Query("department_id")
	
	var datePtr *string
	var deptIDPtr *int
	
	if date != "" {
		datePtr = &date
	}
	
	if deptID != "" {
		var id int
		if _, err := fmt.Sscanf(deptID, "%d", &id); err == nil {
			deptIDPtr = &id
		}
	}
	
	logs, err := c.attendanceRepo.GetAttendanceLogs(datePtr, deptIDPtr)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to fetch attendance logs")
		return
	}
	
	utils.SuccessJSON(ctx, http.StatusOK, "Attendance logs fetched successfully", logs)
}