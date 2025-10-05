package models

import (
	"time"
)

type Attendance struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AttendanceID string     `gorm:"size:100;not null;uniqueIndex" json:"attendance_id"`
	EmployeeID   string     `gorm:"size:50;not null;index" json:"employee_id"`
	ClockIn      time.Time  `gorm:"not null;index" json:"clock_in"`
	ClockInDate  time.Time  `gorm:"type:date;not null;index" json:"clock_in_date"`
	ClockOut     *time.Time `gorm:"index" json:"clock_out"`
	WorkHours    *float64   `gorm:"type:decimal(4,2)" json:"work_hours"`
	Status       string     `gorm:"size:20;default:present" json:"status"` // present, late, half-day, absent
	Notes        string     `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	
	Employee Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}

type AttendanceRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Notes      string `json:"notes"`
}

type ClockOutRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Notes      string `json:"notes"`
}

type AttendanceResponse struct {
	ID           uint              `json:"id"`
	AttendanceID string            `json:"attendance_id"`
	EmployeeID   string            `json:"employee_id"`
	ClockIn      time.Time         `json:"clock_in"`
	ClockInDate  time.Time         `json:"clock_in_date"`
	ClockOut     *time.Time        `json:"clock_out"`
	WorkHours    *float64          `json:"work_hours"`
	Status       string            `json:"status"`
	Notes        string            `json:"notes"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Employee     EmployeeResponse  `json:"employee,omitempty"`
	
	// Punctuality fields
	IsLate       bool   `json:"is_late"`
	LateMinutes  int    `json:"late_minutes,omitempty"`
	IsEarlyLeave bool   `json:"is_early_leave"`
	EarlyMinutes int    `json:"early_minutes,omitempty"`
	Punctuality  string `json:"punctuality"` // "on_time", "late", "early_leave"

	WorkHoursFormatted string `json:"work_hours_formatted,omitempty"`
	ClockInFormatted   string `json:"clock_in_formatted,omitempty"`
	ClockOutFormatted  string `json:"clock_out_formatted,omitempty"`
}

func (a *Attendance) ToResponse() AttendanceResponse {
	return AttendanceResponse{
		ID:           a.ID,
		AttendanceID: a.AttendanceID,
		EmployeeID:   a.EmployeeID,
		ClockIn:      a.ClockIn,
		ClockInDate:  a.ClockInDate,
		ClockOut:     a.ClockOut,
		WorkHours:    a.WorkHours,
		Status:       a.Status,
		Notes:        a.Notes,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		Employee:     a.Employee.ToResponse(),
	}
}