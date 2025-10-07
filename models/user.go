package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password     string     `gorm:"size:255" json:"-"` // Allow empty for setup tokens
	Role         string     `gorm:"size:50;default:employee" json:"role"`
	EmployeeID   *string    `gorm:"size:50;index" json:"employee_id"`
	IsActive     bool       `gorm:"default:false" json:"is_active"` // False until setup complete
	SetupToken   *string    `gorm:"size:64;uniqueIndex" json:"-"`
	TokenExpires *time.Time `json:"-"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Employee *Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
}

type UserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Role       string `json:"role" binding:"required,oneof=admin manager employee"`
	EmployeeID string `json:"employee_id"`
}

type UserSetupRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserSetupResponse struct {
	Message  string `json:"message"`
	LoginURL string `json:"login_url"`
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
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	EmployeeID *string    `json:"employee_id"`
	IsActive   bool       `json:"is_active"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Employee   *EmployeeResponse `json:"employee,omitempty"`
}

type CreateUserResponse struct {
	User       UserResponse `json:"user"`
	SetupToken string       `json:"setup_token,omitempty"`
	Message    string       `json:"message"`
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	if u.Password == "" {
		return nil // Skip hashing if password is empty (for setup tokens)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	if u.Password == "" {
		return false // No password set yet
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate GORM hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword()
}

// BeforeUpdate GORM hook
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Only hash password if it's being updated and not empty
	if tx.Statement.Changed("Password") && u.Password != "" {
		return u.HashPassword()
	}
	return nil
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

// IsSetupTokenValid checks if the setup token is valid and not expired
func (u *User) IsSetupTokenValid() bool {
	if u.SetupToken == nil || u.TokenExpires == nil {
		return false
	}
	return u.TokenExpires.After(time.Now())
}

// ActivateUser activates the user account and clears setup token
func (u *User) ActivateUser() {
	u.IsActive = true
	u.SetupToken = nil
	u.TokenExpires = nil
	u.UpdatedAt = time.Now()
}

// SetSetupToken sets a new setup token with expiration
func (u *User) SetSetupToken(token string, expiresIn time.Duration) {
	u.SetupToken = &token
	expiry := time.Now().Add(expiresIn)
	u.TokenExpires = &expiry
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// CanLogin checks if user can login (active and has password)
func (u *User) CanLogin() bool {
	return u.IsActive && u.Password != ""
}