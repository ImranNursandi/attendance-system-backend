package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Username   string     `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email      string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password   string     `gorm:"size:255;not null" json:"-"`
	Role       string     `gorm:"size:50;default:employee" json:"role"`
	EmployeeID *string    `gorm:"size:50;index" json:"employee_id"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Employee *Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}

type UserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Role       string `json:"role" binding:"required,oneof=admin manager employee"`
	EmployeeID string `json:"employee_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
}

type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	EmployeeID *string   `json:"employee_id"`
	IsActive   bool      `json:"is_active"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Employee   *EmployeeResponse `json:"employee,omitempty"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate GORM hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword()
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	userResp := UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Role:       u.Role,
		EmployeeID: u.EmployeeID,
		IsActive:   u.IsActive,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}

	if u.Employee != nil {
		employeeResp := u.Employee.ToResponse()
		userResp.Employee = &employeeResp
	}

	return userResp
}