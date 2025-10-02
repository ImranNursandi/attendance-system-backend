package models

import (
	"time"
)

type Employee struct {
	ID           string    `gorm:"primaryKey;size:50" json:"id"`
	DepartmentID int       `gorm:"not null" json:"department_id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Address      string    `gorm:"type:text" json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

type Department struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string    `gorm:"size:255;not null" json:"name"`
	MaxClockIn        string    `gorm:"size:5;default:'09:00'" json:"max_clock_in"`    // Format: HH:MM
	MaxClockOut       string    `gorm:"size:5;default:'17:00'" json:"max_clock_out"`   // Format: HH:MM
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}