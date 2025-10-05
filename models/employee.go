package models

import (
	"time"
)

type Employee struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    EmployeeID   string    `gorm:"uniqueIndex;size:50;not null" json:"employee_id"`
    DepartmentID uint      `gorm:"not null;index" json:"department_id"`
    Name         string    `gorm:"size:255;not null" json:"name"`
    Phone        string    `gorm:"size:20" json:"phone"`
    Address      string    `gorm:"type:text;not null" json:"address"`
    Position     string    `gorm:"size:100" json:"position"`
    Status       string    `gorm:"size:20;default:active" json:"status"` // active, inactive, suspended
    JoinDate     time.Time `json:"join_date"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    
    Department   Department   `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
    Attendances  []Attendance `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"attendances,omitempty"`
}

type EmployeeWithUserResponse struct {
	Employee  EmployeeResponse `json:"employee"`
	UserEmail string           `json:"user_email"`
}

type EmployeeRequest struct {
	EmployeeID   string    `json:"employee_id"`
	DepartmentID uint      `json:"department_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address" binding:"required"`
	Position     string    `json:"position"`
	Status       string    `json:"status"`
	JoinDate     time.Time `json:"join_date"`
}

type EmployeeResponse struct {
	ID           uint      `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	DepartmentID uint      `json:"department_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Position     string    `json:"position"`
	Status       string    `json:"status"`
	JoinDate     time.Time `json:"join_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Department   DepartmentResponse `json:"department,omitempty"`
}

func (e *Employee) ToResponse() EmployeeResponse {
	return EmployeeResponse{
		ID:           e.ID,
		EmployeeID:   e.EmployeeID,
		DepartmentID: e.DepartmentID,
		Name:         e.Name,
		Phone:        e.Phone,
		Address:      e.Address,
		Position:     e.Position,
		Status:       e.Status,
		JoinDate:     e.JoinDate,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		Department:   e.Department.ToResponse(),
	}
}