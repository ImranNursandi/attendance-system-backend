package services

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/utils"
	"fmt"
	"time"
)

type AttendanceService struct {
	attendanceRepo *repositories.AttendanceRepository
	employeeRepo   *repositories.EmployeeRepository
}

func NewAttendanceService() *AttendanceService {
	return &AttendanceService{
		attendanceRepo: repositories.NewAttendanceRepository(),
		employeeRepo:   repositories.NewEmployeeRepository(),
	}
}

func (s *AttendanceService) ClockIn(req models.AttendanceRequest) (*models.Attendance, error) {
	// Check if employee exists
	employee, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, utils.NewNotFoundError("employee not found")
	}

	// Check if employee is active
	if employee.Status != "active" {
		return nil, utils.NewBadRequestError("employee is not active")
	}

	// Check if already clocked in today
	todayAttendance, _ := s.attendanceRepo.FindTodayAttendance(req.EmployeeID)
	if todayAttendance != nil && todayAttendance.ID > 0 {
		return nil, utils.NewConflictError("already clocked in today")
	}

	now := time.Now()
	attendance := &models.Attendance{
		EmployeeID: req.EmployeeID,
		ClockIn:    now,
		ClockInDate: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		Notes:      req.Notes,
		Status:     "present",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Check if clock in is on time
	department := employee.Department
	isLate, lateMinutes := utils.CheckClockInLate(now, department.MaxClockIn, department.LateTolerance)

	if isLate {
		attendance.Status = "late"
		attendance.Notes += fmt.Sprintf(" (Late by %d minutes)", lateMinutes)
	}

	// Create attendance record
	if err := s.attendanceRepo.CreateAttendance(attendance); err != nil {
		return nil, err
	}

	// Create attendance history for clock in
	historyIn := &models.AttendanceHistory{
		EmployeeID:     req.EmployeeID,
		DateAttendance: now,
		AttendanceType: 1, // Clock In
		Description:    fmt.Sprintf("Clock In - %s", map[bool]string{true: "Late", false: "On Time"}[isLate]),
		NewValue:       now.Format("2006-01-02 15:04:05"),
		ChangedBy:      "system",
		Reason:         "Regular clock in",
		CreatedAt:      now,
	}

	if err := s.attendanceRepo.CreateAttendanceHistory(historyIn); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) ClockOut(req models.ClockOutRequest) (*models.Attendance, error) {
	// Find today's attendance
	attendance, err := s.attendanceRepo.FindTodayAttendance(req.EmployeeID)
	if err != nil {
		return nil, err
	}
	if attendance == nil {
		return nil, utils.NewNotFoundError("no clock in record found for today")
	}

	if attendance.ClockOut != nil {
		return nil, utils.NewConflictError("already clocked out today")
	}

	now := time.Now()
	attendance.ClockOut = &now
	attendance.UpdatedAt = now

	if req.Notes != "" {
		if attendance.Notes != "" {
			attendance.Notes += " | " + req.Notes
		} else {
			attendance.Notes = req.Notes
		}
	}

	// Calculate work hours
	workHours := now.Sub(attendance.ClockIn).Hours()
	attendance.WorkHours = &workHours

	// Check if clock out is early
	employee, _ := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
	isEarlyLeave, earlyMinutes := false, 0
	if employee != nil {
		isEarlyLeave, earlyMinutes = utils.CheckClockOutEarly(now, employee.Department.MaxClockOut, employee.Department.EarlyLeavePenalty)
		if isEarlyLeave {
			attendance.Notes += fmt.Sprintf(" (Left early by %d minutes)", earlyMinutes)
		}
	}

	// Update attendance record
	if err := s.attendanceRepo.UpdateAttendance(attendance); err != nil {
		return nil, err
	}

	// Create attendance history for clock out
	historyOut := &models.AttendanceHistory{
		EmployeeID:     req.EmployeeID,
		DateAttendance: now,
		AttendanceType: 2, // Clock Out
		Description:    fmt.Sprintf("Clock Out - %s", map[bool]string{true: "Early Leave", false: "On Time"}[isEarlyLeave]),
		NewValue:       now.Format("2006-01-02 15:04:05"),
		ChangedBy:      "system",
		Reason:         "Regular clock out",
		CreatedAt:      now,
	}

	if err := s.attendanceRepo.CreateAttendanceHistory(historyOut); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) GetAttendanceLogs(startDate, endDate string, departmentID uint, employeeID string, page, limit int) ([]models.Attendance, *repositories.Pagination, error) {
	return s.attendanceRepo.GetAttendanceLogs(startDate, endDate, departmentID, employeeID, page, limit)
}

func (s *AttendanceService) GetEmployeeAttendance(employeeID string, startDate, endDate string) ([]models.Attendance, error) {
	return s.attendanceRepo.GetEmployeeAttendance(employeeID, startDate, endDate)
}

func (s *AttendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	return s.attendanceRepo.GetAttendanceByID(id)
}

func (s *AttendanceService) GetAttendanceStats(employeeID string, month, year int) (map[string]interface{}, error) {
	if month == 0 {
		month = int(time.Now().Month())
	}
	if year == 0 {
		year = time.Now().Year()
	}
	return s.attendanceRepo.GetAttendanceStats(employeeID, month, year)
}