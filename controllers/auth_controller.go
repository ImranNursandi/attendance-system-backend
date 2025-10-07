package controllers

import (
	"attendance-system/models"
	"attendance-system/services"
	"attendance-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// SetupAccount godoc
// @Summary Setup user account
// @Description Activate user account and set password using setup token
// @Tags auth
// @Accept json
// @Produce json
// @Param setupData body models.UserSetupRequest true "Account setup data"
// @Success 200 {object} utils.Response{data=models.UserSetupResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 410 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/setup-account [post]
func (c *AuthController) SetupAccount(ctx *gin.Context) {
	var req models.UserSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	response, err := c.authService.SetupAccount(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Account setup completed successfully", response)
}

// VerifySetupToken godoc
// @Summary Verify setup token
// @Description Check if setup token is valid and not expired
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Setup token"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 410 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/verify-setup-token [get]
func (c *AuthController) VerifySetupToken(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Setup token is required")
		return
	}

	user, err := c.authService.VerifySetupToken(token)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	response := map[string]interface{}{
		"valid":      true,
		"email":      user.Email,
		"username":   user.Username,
		"expires_at": user.TokenExpires,
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Setup token is valid", response)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset instructions to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body map[string]string true "User email"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := c.authService.ForgotPassword(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		utils.SuccessJSON(ctx, http.StatusOK, "If the email exists, password reset instructions have been sent", nil)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "If the email exists, password reset instructions have been sent", nil)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param resetData body models.ResetPasswordRequest true "Password reset data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 410 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := c.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Password reset successfully", nil)
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token (client-side)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// We could implement token blacklisting if needed
	utils.SuccessJSON(ctx, http.StatusOK, "Logged out successfully", nil)
}

// Existing methods remain the same...

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=models.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	response, err := c.authService.Login(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Login successful", response)
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User registration data"
// @Success 201 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	response, err := c.authService.Register(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "User registered successfully", response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/profile [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	response, err := c.authService.GetUserProfile(uint(userIDUint))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Profile retrieved successfully", response)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body map[string]interface{} true "Profile update data"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/profile [put]
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req map[string]interface{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	response, err := c.authService.UpdateProfile(uint(userIDUint), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Profile updated successfully", response)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body map[string]string true "Password change data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/change-password [put]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err = c.authService.ChangePassword(uint(userIDUint), req.CurrentPassword, req.NewPassword)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Password changed successfully", nil)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh expired JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// This would typically use a refresh token
	// For now, we'll return an error indicating this feature is not implemented
	utils.ErrorJSON(ctx, http.StatusNotImplemented, "Refresh token endpoint not implemented")
}