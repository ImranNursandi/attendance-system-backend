package models

import (
	"time"
)

type AttendanceHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     string    `gorm:"size:50;not null;index" json:"employee_id"`
	AttendanceID   string    `gorm:"size:100;not null;index" json:"attendance_id"`
	DateAttendance time.Time `gorm:"not null;index" json:"date_attendance"`
	AttendanceType int8      `gorm:"type:tinyint;not null" json:"attendance_type"` // 1: Clock In, 2: Clock Out, 3: Adjustment, 4: Correction
	Description    string    `gorm:"type:text" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	
	Employee   Employee   `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
	Attendance Attendance `gorm:"foreignKey:AttendanceID;references:AttendanceID" json:"attendance,omitempty"`
}

type AttendanceHistoryRequest struct {
	EmployeeID     string    `json:"employee_id" binding:"required"`
	AttendanceID   string    `json:"attendance_id" binding:"required"`
	DateAttendance time.Time `json:"date_attendance" binding:"required"`
	AttendanceType int8      `json:"attendance_type" binding:"required"`
	Description    string    `json:"description"`
}

type AttendanceHistoryResponse struct {
	ID             uint      `json:"id"`
	EmployeeID     string    `json:"employee_id"`
	AttendanceID   string    `json:"attendance_id"`
	DateAttendance time.Time `json:"date_attendance"`
	AttendanceType int8      `json:"attendance_type"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Employee       EmployeeResponse   `json:"employee,omitempty"`
	Attendance     AttendanceResponse `json:"attendance,omitempty"`
}

func (h *AttendanceHistory) ToResponse() AttendanceHistoryResponse {
	return AttendanceHistoryResponse{
		ID:             h.ID,
		EmployeeID:     h.EmployeeID,
		AttendanceID:   h.AttendanceID,
		DateAttendance: h.DateAttendance,
		AttendanceType: h.AttendanceType,
		Description:    h.Description,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
		Employee:       h.Employee.ToResponse(),
		Attendance:     h.Attendance.ToResponse(),
	}
}