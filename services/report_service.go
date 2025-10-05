package services

import (
	"attendance-system/repositories"
	"fmt"
	"time"
)

type ReportService struct {
	attendanceRepo *repositories.AttendanceRepository
	employeeRepo   *repositories.EmployeeRepository
	departmentRepo *repositories.DepartmentRepository
}

func NewReportService() *ReportService {
	return &ReportService{
		attendanceRepo: repositories.NewAttendanceRepository(),
		employeeRepo:   repositories.NewEmployeeRepository(),
		departmentRepo: repositories.NewDepartmentRepository(),
	}
}

type AttendanceReport struct {
	EmployeeID    string    `json:"employee_id"`
	EmployeeName  string    `json:"employee_name"`
	Department    string    `json:"department"`
	Date          string    `json:"date"`
	ClockIn       time.Time `json:"clock_in"`
	ClockOut      time.Time `json:"clock_out"`
	WorkHours     float64   `json:"work_hours"`
	Status        string    `json:"status"`
	LateMinutes   int       `json:"late_minutes"`
	EarlyMinutes  int       `json:"early_minutes"`
}

type SummaryReport struct {
	Period           string `json:"period"`
	TotalEmployees   int64  `json:"total_employees"`
	TotalPresent     int64  `json:"total_present"`
	TotalLate        int64  `json:"total_late"`
	TotalAbsent      int64  `json:"total_absent"`
	TotalWorkHours   string `json:"total_work_hours"`
	AverageWorkHours string `json:"average_work_hours"`
}

func (s *ReportService) GenerateAttendanceReport(startDate, endDate string, departmentID uint) ([]AttendanceReport, error) {
	var reports []AttendanceReport

	// Get attendances for the period
	attendances, _, err := s.attendanceRepo.GetAttendanceLogs(startDate, endDate, departmentID, "", 1, 10000)
	if err != nil {
		return nil, err
	}

	for _, attendance := range attendances {
		report := AttendanceReport{
			EmployeeID:   attendance.EmployeeID,
			EmployeeName: attendance.Employee.Name,
			Department:   attendance.Employee.Department.Name,
			Date:         attendance.ClockIn.Format("2006-01-02"),
			ClockIn:      attendance.ClockIn,
			Status:       attendance.Status,
		}

		if attendance.ClockOut != nil {
			report.ClockOut = *attendance.ClockOut
			if attendance.WorkHours != nil {
				report.WorkHours = *attendance.WorkHours
			}
		}

		// Calculate late minutes
		maxClockIn, _ := time.Parse("15:04:05", attendance.Employee.Department.MaxClockIn)
		clockInTime := time.Date(attendance.ClockIn.Year(), attendance.ClockIn.Month(), attendance.ClockIn.Day(),
			attendance.ClockIn.Hour(), attendance.ClockIn.Minute(), attendance.ClockIn.Second(), 0, attendance.ClockIn.Location())
		maxClockInTime := time.Date(attendance.ClockIn.Year(), attendance.ClockIn.Month(), attendance.ClockIn.Day(),
			maxClockIn.Hour(), maxClockIn.Minute(), maxClockIn.Second(), 0, attendance.ClockIn.Location())

		if clockInTime.After(maxClockInTime) {
			report.LateMinutes = int(clockInTime.Sub(maxClockInTime).Minutes())
		}

		// Calculate early leave minutes
		if attendance.ClockOut != nil {
			maxClockOut, _ := time.Parse("15:04:05", attendance.Employee.Department.MaxClockOut)
			clockOutTime := time.Date(attendance.ClockOut.Year(), attendance.ClockOut.Month(), attendance.ClockOut.Day(),
				attendance.ClockOut.Hour(), attendance.ClockOut.Minute(), attendance.ClockOut.Second(), 0, attendance.ClockOut.Location())
			maxClockOutTime := time.Date(attendance.ClockOut.Year(), attendance.ClockOut.Month(), attendance.ClockOut.Day(),
				maxClockOut.Hour(), maxClockOut.Minute(), maxClockOut.Second(), 0, attendance.ClockOut.Location())

			if clockOutTime.Before(maxClockOutTime) {
				report.EarlyMinutes = int(maxClockOutTime.Sub(clockOutTime).Minutes())
			}
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (s *ReportService) GenerateSummaryReport(startDate, endDate string, departmentID uint) (*SummaryReport, error) {
	var summary SummaryReport
	summary.Period = fmt.Sprintf("%s to %s", startDate, endDate)

	// Get total employees
	totalEmployees, err := s.employeeRepo.GetActiveEmployeesCount()
	if err != nil {
		return nil, err
	}
	summary.TotalEmployees = totalEmployees

	// Get attendance statistics
	attendances, _, err := s.attendanceRepo.GetAttendanceLogs(startDate, endDate, departmentID, "", 1, 10000)
	if err != nil {
		return nil, err
	}

	var totalPresent, totalLate int64
	var totalWorkHours float64

	for _, attendance := range attendances {
		totalPresent++
		if attendance.Status == "late" {
			totalLate++
		}
		if attendance.WorkHours != nil {
			totalWorkHours += *attendance.WorkHours
		}
	}

	summary.TotalPresent = totalPresent
	summary.TotalLate = totalLate
	summary.TotalAbsent = totalEmployees - totalPresent

	if totalPresent > 0 {
		avgWorkHours := totalWorkHours / float64(totalPresent)
		summary.AverageWorkHours = fmt.Sprintf("%.2f hours", avgWorkHours)
		summary.TotalWorkHours = fmt.Sprintf("%.2f hours", totalWorkHours)
	} else {
		summary.AverageWorkHours = "0 hours"
		summary.TotalWorkHours = "0 hours"
	}

	return &summary, nil
}

func (s *ReportService) GenerateDepartmentReport(departmentID uint, month, year int) (map[string]interface{}, error) {
	if month == 0 {
		month = int(time.Now().Month())
	}
	if year == 0 {
		year = time.Now().Year()
	}

	// startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// endDate := startDate.AddDate(0, 1, -1)

	// Get department employees
	employees, err := s.employeeRepo.FindByDepartment(departmentID)
	if err != nil {
		return nil, err
	}

	var departmentStats struct {
		DepartmentName string                   `json:"department_name"`
		TotalEmployees int                      `json:"total_employees"`
		EmployeeStats  []map[string]interface{} `json:"employee_stats"`
		Summary        map[string]interface{}   `json:"summary"`
	}

	department, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return nil, err
	}

	departmentStats.DepartmentName = department.Name
	departmentStats.TotalEmployees = len(employees)
	departmentStats.EmployeeStats = make([]map[string]interface{}, 0)

	var totalPresent, totalLate, totalAbsent int64
	var totalWorkHours float64

	for _, employee := range employees {
		stats, err := s.attendanceRepo.GetAttendanceStats(employee.EmployeeID, month, year)
		if err != nil {
			continue
		}

		employeeStat := map[string]interface{}{
			"employee_id":   employee.EmployeeID,
			"employee_name": employee.Name,
			"stats":         stats,
		}

		departmentStats.EmployeeStats = append(departmentStats.EmployeeStats, employeeStat)

		totalPresent += stats["total_present"].(int64)
		totalLate += stats["total_late"].(int64)
		totalAbsent += stats["total_absent"].(int64)
		if avgHours, ok := stats["avg_work_hours"].(float64); ok {
			totalWorkHours += avgHours
		}
	}

	// Calculate attendance rate safely
	attendanceRate := 0.0
	if totalPresent+totalAbsent > 0 {
		attendanceRate = float64(totalPresent) / float64(totalPresent+totalAbsent) * 100
	}

	// Calculate average work hours safely
	avgWorkHours := 0.0
	if len(employees) > 0 {
		avgWorkHours = totalWorkHours / float64(len(employees))
	}

	departmentStats.Summary = map[string]interface{}{
		"total_present":      totalPresent,
		"total_late":         totalLate,
		"total_absent":       totalAbsent,
		"attendance_rate":    attendanceRate,
		"average_work_hours": avgWorkHours,
	}

	return map[string]interface{}{
		"department_report": departmentStats,
		"period":            fmt.Sprintf("%d-%02d", year, month),
	}, nil
}