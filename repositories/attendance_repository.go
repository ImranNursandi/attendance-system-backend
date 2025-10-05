package repositories

import (
	"attendance-system/models"
	"time"
)

type AttendanceRepository struct {
	BaseRepository
}

func NewAttendanceRepository() *AttendanceRepository {
	return &AttendanceRepository{
		BaseRepository: *NewBaseRepository(),
	}
}

func (r *AttendanceRepository) CreateAttendance(attendance *models.Attendance) error {
	if err := r.DB.Create(attendance).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *AttendanceRepository) FindTodayAttendance(employeeID string) (*models.Attendance, error) {
	var attendance models.Attendance
	today := time.Now().Format("2006-01-02")
	
	err := r.DB.Preload("Employee.Department").
		Where("employee_id = ? AND DATE(clock_in) = ?", employeeID, today).
		First(&attendance).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return &attendance, nil
}

func (r *AttendanceRepository) UpdateAttendance(attendance *models.Attendance) error {
	if err := r.DB.Save(attendance).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *AttendanceRepository) CreateAttendanceHistory(history *models.AttendanceHistory) error {
	if err := r.DB.Create(history).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *AttendanceRepository) GetAttendanceLogs(startDate, endDate string, departmentID uint, employeeID string, page, limit int) ([]models.Attendance, *Pagination, error) {
	var attendances []models.Attendance
	
	query := r.DB.Preload("Employee.Department")
	
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(clock_in) BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		query = query.Where("DATE(clock_in) >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("DATE(clock_in) <= ?", endDate)
	}
	
	if departmentID > 0 {
		query = query.Joins("JOIN employees ON attendances.employee_id = employees.employee_id").
			Where("employees.department_id = ?", departmentID)
	}
	
	if employeeID != "" {
		query = query.Where("attendances.employee_id = ?", employeeID)
	}
	
	pagination, err := r.Paginate(query.Order("clock_in DESC"), page, limit, &attendances)
	if err != nil {
		return nil, nil, r.HandleError(err)
	}
	
	return attendances, pagination, nil
}

func (r *AttendanceRepository) GetEmployeeAttendance(employeeID string, startDate, endDate string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	
	query := r.DB.Preload("Employee.Department").
		Where("employee_id = ?", employeeID)
	
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(clock_in) BETWEEN ? AND ?", startDate, endDate)
	}
	
	err := query.Order("clock_in DESC").Find(&attendances).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	
	return attendances, nil
}

func (r *AttendanceRepository) GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.DB.Preload("Employee.Department").First(&attendance, id).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return &attendance, nil
}

func (r *AttendanceRepository) GetAttendanceStats(employeeID string, month, year int) (map[string]interface{}, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	
	var stats struct {
		TotalPresent  int64
		TotalLate     int64
		TotalAbsent   int64
		TotalWorkDays int64
		AvgWorkHours  float64
	}
	
	// Count present days
	r.DB.Model(&models.Attendance{}).
		Where("employee_id = ? AND DATE(clock_in) BETWEEN ? AND ?", employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Count(&stats.TotalPresent)
	
	// Count late days
	r.DB.Model(&models.Attendance{}).
		Where("employee_id = ? AND DATE(clock_in) BETWEEN ? AND ? AND status = ?", 
			employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), "late").
		Count(&stats.TotalLate)
	
	// Calculate total work days in month
	stats.TotalWorkDays = int64(endDate.Day())
	
	// Calculate absent days
	stats.TotalAbsent = stats.TotalWorkDays - stats.TotalPresent
	
	// Calculate average work hours
	row := r.DB.Table("attendances").
		Where("employee_id = ? AND DATE(clock_in) BETWEEN ? AND ? AND clock_out IS NOT NULL", 
			employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("AVG(TIMESTAMPDIFF(HOUR, clock_in, clock_out))").
		Row()
	row.Scan(&stats.AvgWorkHours)
	
	return map[string]interface{}{
		"total_present":  stats.TotalPresent,
		"total_late":     stats.TotalLate,
		"total_absent":   stats.TotalAbsent,
		"total_work_days": stats.TotalWorkDays,
		"avg_work_hours": stats.AvgWorkHours,
	}, nil
}