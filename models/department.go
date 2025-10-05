package models

import (
	"time"
)

// In department.go
type Department struct {
    ID                uint      `gorm:"primaryKey" json:"id"`
    Name              string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
    Description       string    `gorm:"type:text" json:"description"`
    MaxClockIn        string    `gorm:"size:8;not null" json:"max_clock_in"`        // Format: HH:MM:SS
    MaxClockOut       string    `gorm:"size:8;not null" json:"max_clock_out"`       // Format: HH:MM:SS
    LateTolerance     int       `gorm:"default:0" json:"late_tolerance"`            // Tolerance in minutes
    EarlyLeavePenalty int       `gorm:"default:0" json:"early_leave_penalty"`       // Penalty threshold in minutes
    Status            string    `gorm:"size:20;default:active" json:"status"`       // active, inactive
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
    
    Employees []Employee `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`
}

type DepartmentRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	MaxClockIn        string `json:"max_clock_in" binding:"required"`
	MaxClockOut       string `json:"max_clock_out" binding:"required"`
	LateTolerance     int    `json:"late_tolerance"`
	EarlyLeavePenalty int    `json:"early_leave_penalty"`
	Status            string `json:"status"`
}

type DepartmentResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	MaxClockIn        string    `json:"max_clock_in"`
	MaxClockOut       string    `json:"max_clock_out"`
	LateTolerance     int       `json:"late_tolerance"`
	EarlyLeavePenalty int       `json:"early_leave_penalty"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	EmployeeCount     int       `json:"employee_count,omitempty"`
}

func (d *Department) ToResponse() DepartmentResponse {
	return DepartmentResponse{
		ID:                d.ID,
		Name:              d.Name,
		Description:       d.Description,
		MaxClockIn:        d.MaxClockIn,
		MaxClockOut:       d.MaxClockOut,
		LateTolerance:     d.LateTolerance,
		EarlyLeavePenalty: d.EarlyLeavePenalty,
		Status:            d.Status,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}
}