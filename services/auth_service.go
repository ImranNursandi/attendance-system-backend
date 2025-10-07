package services

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo     *repositories.UserRepository
	employeeRepo *repositories.EmployeeRepository
	emailService *ResendEmailService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:     repositories.NewUserRepository(),
		employeeRepo: repositories.NewEmployeeRepository(),
		emailService: NewEmailService(),
	}
}

// SetupAccount activates user account and sets password using setup token
func (s *AuthService) SetupAccount(req models.UserSetupRequest) (*models.UserSetupResponse, error) {
	// Find user by setup token
	user, err := s.userRepo.FindBySetupToken(req.Token)
	if err != nil {
		return nil, utils.NewNotFoundError("invalid or expired setup token")
	}

	// Check if token is still valid
	if !user.IsSetupTokenValid() {
		return nil, utils.NewBadRequestError("setup token has expired")
	}

	// Check if user is already active
	if user.IsActive {
		return nil, utils.NewBadRequestError("account is already active")
	}

	// Update user with new password and activate account
	user.Password = req.NewPassword
	user.ActivateUser()
	user.UpdatedAt = time.Now()

	if err := user.HashPassword(); err != nil {
		return nil, utils.NewInternalServerError("failed to hash password")
	}

	// Save the updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, utils.NewInternalServerError("failed to setup account")
	}

	// Send welcome email
	if user.Email != "" {
		employeeName := user.Username
		if user.Employee != nil {
			employeeName = user.Employee.Name
		}
		go s.emailService.SendWelcomeEmail(user.Email, employeeName)
	}

	response := &models.UserSetupResponse{
		Message:  "Account setup completed successfully. You can now log in.",
		LoginURL: s.getFrontendURL() + "/login",
	}

	return response, nil
}

// VerifySetupToken checks if setup token is valid and returns user info
func (s *AuthService) VerifySetupToken(token string) (*models.User, error) {
	if token == "" {
		return nil, utils.NewBadRequestError("setup token is required")
	}

	user, err := s.userRepo.FindBySetupToken(token)
	if err != nil {
		return nil, utils.NewNotFoundError("invalid setup token")
	}

	// Check if token is expired
	if !user.IsSetupTokenValid() {
		return nil, utils.NewBadRequestError("setup token has expired")
	}

	// Check if user is already active
	if user.IsActive {
		return nil, utils.NewBadRequestError("account is already active")
	}

	return user, nil
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil
	}

	// Check if user is active
	if !user.IsActive {
		// Don't reveal that the account is inactive
		return nil
	}

	// Generate reset token
	resetToken, err := models.GenerateSecureToken(32)
	if err != nil {
		return utils.NewInternalServerError("failed to generate reset token")
	}

	// Set reset token (expires in 1 hour)
	resetExpiry := 1 * time.Hour
	user.SetupToken = &resetToken
	expiry := time.Now().Add(resetExpiry)
	user.TokenExpires = &expiry
	user.UpdatedAt = time.Now()

	// Save user with reset token
	if err := s.userRepo.Update(user); err != nil {
		return utils.NewInternalServerError("failed to set reset token")
	}

	// Send password reset email
	go s.emailService.SendPasswordResetEmail(user.Email, resetToken)

	return nil
}

// ResetPassword resets user password using reset token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Find user by reset token
	user, err := s.userRepo.FindBySetupToken(token)
	if err != nil {
		return utils.NewNotFoundError("invalid or expired reset token")
	}

	// Check if token is still valid
	if user.TokenExpires == nil || user.TokenExpires.Before(time.Now()) {
		return utils.NewBadRequestError("reset token has expired")
	}

	// Update password and clear reset token
	user.Password = newPassword
	user.SetupToken = nil
	user.TokenExpires = nil
	user.UpdatedAt = time.Now()

	// Save the updated user
	if err := s.userRepo.Update(user); err != nil {
		return utils.NewInternalServerError("failed to reset password")
	}

	return nil
}

// Enhanced Login method to handle inactive accounts with setup tokens
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	var err error

	// Try finding by username first
	user, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		// If not found by username, try by email
		user, err = s.userRepo.FindByEmail(req.Username)
		if err != nil {
			return nil, utils.NewUnauthorizedError("Invalid Credentials")
		}
	}

	// Check if user is active
	if !user.IsActive {
		// If user has a valid setup token, guide them to setup
		if user.IsSetupTokenValid() {
			return nil, utils.NewUnauthorizedError("Account setup required. Please check your email for setup instructions.")
		}
		return nil, utils.NewUnauthorizedError("Account Is Deactivated")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, utils.NewUnauthorizedError("Invalid Credentials")
	}

	// Generate JWT token
	accessToken, expiresIn, err := s.generateJWTToken(user)
	if err != nil {
		return nil, utils.NewInternalServerError("failed to generate token")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log but don't fail the login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Create response
	response := &models.LoginResponse{
		User:        user.ToResponse(),
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}

	return response, nil
}

// Helper method to get frontend URL
func (s *AuthService) getFrontendURL() string {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	return frontendURL
}

func (s *AuthService) Register(req models.UserRequest) (*models.UserResponse, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, utils.NewConflictError("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, utils.NewConflictError("email already exists")
	}

	// If employee_id is provided, verify it exists
	var employeeID *string
	if req.EmployeeID != "" {
		employee, err := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
		if err != nil {
			return nil, utils.NewBadRequestError("employee not found")
		}
		employeeID = &employee.EmployeeID
	}

	// Create user
	user := &models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		EmployeeID: employeeID,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Reload user to get relations
	createdUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	response := createdUser.ToResponse()
	return &response, nil
}

func (s *AuthService) generateJWTToken(user *models.User) (string, int64, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-key-change-in-production"
	}

	// Get token expiry from environment (default 24 hours)
	tokenExpiryStr := os.Getenv("JWT_EXPIRY")
	if tokenExpiryStr == "" {
		tokenExpiryStr = "24h"
	}

	// Parse duration
	var expiryDuration time.Duration
	var err error
	if tokenExpiryStr != "" {
		expiryDuration, err = time.ParseDuration(tokenExpiryStr)
		if err != nil {
			expiryDuration = 24 * time.Hour // Default to 24 hours
		}
	} else {
		expiryDuration = 24 * time.Hour
	}

	expiryTime := time.Now().Add(expiryDuration)
	expiresIn := expiryTime.Unix()

	// Create claims
	claims := utils.JWTClaims{
		UserID:   strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "attendance-system",
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return nil, utils.NewUnauthorizedError("invalid token")
	}
	return claims, nil
}

func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return utils.NewNotFoundError("user not found")
	}

	// Verify current password
	if !user.CheckPassword(currentPassword) {
		return utils.NewBadRequestError("current password is incorrect")
	}

	// Update password
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return utils.NewInternalServerError("failed to hash password")
	}

	if err := s.userRepo.Update(user); err != nil {
		return utils.NewInternalServerError("failed to update password")
	}

	return nil
}

func (s *AuthService) GetUserProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	response := user.ToResponse()
	return &response, nil
}

func (s *AuthService) UpdateProfile(userID uint, req map[string]interface{}) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	if email, ok := req["email"].(string); ok && email != "" {
		// Check if email is already taken by another user
		existingUser, _ := s.userRepo.FindByEmail(email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, utils.NewConflictError("email already taken")
		}
		user.Email = email
	}

	if username, ok := req["username"].(string); ok && username != "" {
		// Check if username is already taken by another user
		existingUser, _ := s.userRepo.FindByUsername(username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, utils.NewConflictError("username already taken")
		}
		user.Username = username
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}