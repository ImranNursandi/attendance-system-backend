package models

import (
	"time"
)

type AttendanceHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	EmployeeID      string    `gorm:"size:50;not null;index" json:"employee_id"`
	DateAttendance  time.Time `gorm:"not null;index" json:"date_attendance"`
	AttendanceType  int8      `gorm:"type:tinyint;not null" json:"attendance_type"` // 1: Clock In, 2: Clock Out, 3: Adjustment, 4: Correction
	Description     string    `gorm:"type:text" json:"description"`
	PreviousValue   string    `gorm:"type:text" json:"previous_value"`
	NewValue        string    `gorm:"type:text" json:"new_value"`
	ChangedBy       string    `gorm:"size:50" json:"changed_by"`
	Reason          string    `gorm:"type:text" json:"reason"`
	CreatedAt       time.Time `json:"created_at"`
	
	Employee Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}

type AttendanceHistoryRequest struct {
	EmployeeID     string    `json:"employee_id" binding:"required"`
	DateAttendance time.Time `json:"date_attendance" binding:"required"`
	AttendanceType int8      `json:"attendance_type" binding:"required"`
	Description    string    `json:"description"`
	PreviousValue  string    `json:"previous_value"`
	NewValue       string    `json:"new_value"`
	ChangedBy      string    `json:"changed_by"`
	Reason         string    `json:"reason"`
}

type AttendanceHistoryResponse struct {
	ID              uint      `json:"id"`
	EmployeeID      string    `json:"employee_id"`
	DateAttendance  time.Time `json:"date_attendance"`
	AttendanceType  int8      `json:"attendance_type"`
	Description     string    `json:"description"`
	PreviousValue   string    `json:"previous_value"`
	NewValue        string    `json:"new_value"`
	ChangedBy       string    `json:"changed_by"`
	Reason          string    `json:"reason"`
	CreatedAt       time.Time `json:"created_at"`
	Employee        EmployeeResponse `json:"employee,omitempty"`
}

func (h *AttendanceHistory) ToResponse() AttendanceHistoryResponse {
	return AttendanceHistoryResponse{
		ID:             h.ID,
		EmployeeID:     h.EmployeeID,
		DateAttendance: h.DateAttendance,
		AttendanceType: h.AttendanceType,
		Description:    h.Description,
		PreviousValue:  h.PreviousValue,
		NewValue:       h.NewValue,
		ChangedBy:      h.ChangedBy,
		Reason:         h.Reason,
		CreatedAt:      h.CreatedAt,
		Employee:       h.Employee.ToResponse(),
	}
}