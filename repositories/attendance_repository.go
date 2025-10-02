package repositories

import (
	"attendance-system/models"
	"time"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) ClockIn(attendance *models.Attendance, history *models.AttendanceHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(attendance).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}

func (r *AttendanceRepository) ClockOut(attendanceID string, clockOut time.Time, history *models.AttendanceHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Attendance{}).Where("id = ?", attendanceID).Update("clock_out", clockOut).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}

func (r *AttendanceRepository) GetAttendanceLogs(date *string, departmentID *int) ([]models.AttendanceHistory, error) {
	var histories []models.AttendanceHistory
	
	query := r.db.Preload("Employee.Department")
	
	if date != nil {
		query = query.Where("DATE(date_attendance) = ?", *date)
	}
	
	if departmentID != nil {
		query = query.Joins("JOIN employees ON employees.id = attendance_histories.employee_id").
			Where("employees.department_id = ?", *departmentID)
	}
	
	err := query.Order("date_attendance DESC").Find(&histories).Error
	return histories, err
}

func (r *AttendanceRepository) GetAttendanceByEmployeeID(employeeID string, date time.Time) (models.Attendance, error) {
	var attendance models.Attendance
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	err := r.db.Where("employee_id = ? AND clock_in >= ? AND clock_in < ?", employeeID, startOfDay, endOfDay).First(&attendance).Error
	return attendance, err
}