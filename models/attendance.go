package models

import (
	"time"
)

type Attendance struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	EmployeeID string    `gorm:"size:50;not null" json:"employee_id"`
	ClockIn    time.Time `json:"clock_in"`
	ClockOut   *time.Time `json:"clock_out"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	Employee   Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

type AttendanceHistory struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	EmployeeID      string    `gorm:"size:50;not null" json:"employee_id"`
	AttendanceID    string    `gorm:"size:100" json:"attendance_id"`
	DateAttendance  time.Time `json:"date_attendance"`
	AttendanceType  int8      `gorm:"type:tinyint(1)" json:"attendance_type"` // 1: Masuk, 2: Keluar, 3: Lainnya
	Description     string    `gorm:"type:text" json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	Employee        Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}