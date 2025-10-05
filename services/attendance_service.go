// services/attendance_service.go
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
	
	// Generate unique attendance ID
	attendanceID := fmt.Sprintf("ATT-%s-%d", req.EmployeeID, now.Unix())
	
	attendance := &models.Attendance{
		AttendanceID: attendanceID, // ← ADDED
		EmployeeID:   req.EmployeeID,
		ClockIn:      now,
		ClockInDate:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		Notes:        req.Notes,
		Status:       "present",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Check if clock in is on time
	department := employee.Department
	isLate, lateMinutes := utils.CheckClockInLate(now, department.MaxClockIn, department.LateTolerance)

	if isLate {
		attendance.Status = "late"
		if attendance.Notes != "" {
			attendance.Notes += fmt.Sprintf(" | Late by %d minutes", lateMinutes)
		} else {
			attendance.Notes = fmt.Sprintf("Late by %d minutes", lateMinutes)
		}
	}

	// Create attendance record
	if err := s.attendanceRepo.CreateAttendance(attendance); err != nil {
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
	if employee != nil {
		isEarlyLeave, earlyMinutes := utils.CheckClockOutEarly(now, employee.Department.MaxClockOut, employee.Department.EarlyLeavePenalty)
		if isEarlyLeave {
			if attendance.Notes != "" {
				attendance.Notes += fmt.Sprintf(" | Left early by %d minutes", earlyMinutes)
			} else {
				attendance.Notes = fmt.Sprintf("Left early by %d minutes", earlyMinutes)
			}
		}
	}

	// Update attendance record
	if err := s.attendanceRepo.UpdateAttendance(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// Enhanced method with punctuality data
func (s *AttendanceService) GetAttendanceLogs(startDate, endDate string, departmentID uint, employeeID string, page, limit int) ([]models.AttendanceResponse, *repositories.Pagination, error) {
	attendances, pagination, err := s.attendanceRepo.GetAttendanceLogs(startDate, endDate, departmentID, employeeID, page, limit)
	if err != nil {
		return nil, nil, err
	}

	// Convert to response with punctuality data
	var responses []models.AttendanceResponse
	for _, attendance := range attendances {
		response := attendance.ToResponse()
		
		// Calculate punctuality
		if attendance.Employee.ID > 0 {
			dept := attendance.Employee.Department
			
			// Check clock-in punctuality
			isLate, lateMinutes := utils.CheckClockInLate(attendance.ClockIn, dept.MaxClockIn, dept.LateTolerance)
			
			// Check clock-out punctuality
			isEarlyLeave, earlyMinutes := false, 0
			if attendance.ClockOut != nil {
				isEarlyLeave, earlyMinutes = utils.CheckClockOutEarly(*attendance.ClockOut, dept.MaxClockOut, dept.EarlyLeavePenalty)
			}
			
			// Set punctuality fields
			response.IsLate = isLate
			response.LateMinutes = lateMinutes
			response.IsEarlyLeave = isEarlyLeave
			response.EarlyMinutes = earlyMinutes
			
			// Overall punctuality status
			if isLate {
				response.Punctuality = "late"
			} else if isEarlyLeave {
				response.Punctuality = "early_leave"
			} else {
				response.Punctuality = "on_time"
			}
		}
		
		responses = append(responses, response)
	}

	return responses, pagination, nil
}

func (s *AttendanceService) GetEmployeeAttendance(employeeID string, startDate, endDate string) ([]models.Attendance, error) {
	return s.attendanceRepo.GetEmployeeAttendance(employeeID, startDate, endDate)
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

func (s *AttendanceService) CalculateAttendancePunctuality(attendance *models.Attendance) *models.AttendanceResponse {
	response := attendance.ToResponse()
	
	if attendance.Employee.ID > 0 {
		dept := attendance.Employee.Department
		
		isLate, lateMinutes, isEarlyLeave, earlyMinutes, punctuality := utils.CalculatePunctualityStatus(
			attendance.ClockIn,
			attendance.ClockOut,
			dept.MaxClockIn,
			dept.MaxClockOut,
			dept.LateTolerance,
			dept.EarlyLeavePenalty,
		)
		
		response.IsLate = isLate
		response.LateMinutes = lateMinutes
		response.IsEarlyLeave = isEarlyLeave
		response.EarlyMinutes = earlyMinutes
		response.Punctuality = punctuality
	}
	
	return &response
}

func (s *AttendanceService) GetEmployeeAttendanceWithPunctuality(employeeID string, startDate, endDate string) ([]models.AttendanceResponse, error) {
	attendances, err := s.attendanceRepo.GetEmployeeAttendance(employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var responses []models.AttendanceResponse
	for _, attendance := range attendances {
		response := s.CalculateAttendancePunctuality(&attendance)
		responses = append(responses, *response)
	}

	return responses, nil
}