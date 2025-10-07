package controllers

import (
	"attendance-system/models"
	"attendance-system/services"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SetupController struct {
	employeeService *services.EmployeeService
}

func NewSetupController() *SetupController {
	return &SetupController{
		employeeService: services.NewEmployeeService(),
	}
}

// CompleteAccountSetup handles account activation and password setting
func (sc *SetupController) CompleteAccountSetup(c *gin.Context) {
	var req models.UserSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.employeeService.CompleteAccountSetup(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := models.UserSetupResponse{
		Message:  "Account setup completed successfully. You can now login.",
		LoginURL: "https://yourapp.com/login",
	}

	c.JSON(http.StatusOK, response)
}

// GetSetupToken allows managers to retrieve setup token if email fails
func (sc *SetupController) GetSetupToken(c *gin.Context) {
	employeeID := c.Param("employeeId")
	
	token, err := sc.employeeService.GetEmployeeSetupToken(employeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setupURL := fmt.Sprintf("https://yourapp.com/setup-account?token=%s", token)
	
	c.JSON(http.StatusOK, gin.H{
		"setup_token": token,
		"setup_url":   setupURL,
		"expires_in":  "7 days",
		"instructions": "Provide this URL to the employee to set up their account",
	})
}

// func (s *SetupController) ResendSetupEmail(employeeID string) error {
// 	user, err := s.userRepo.FindByEmployeeID(employeeID)
// 	if err != nil {
// 		return utils.NewNotFoundError("Employee user account not found")
// 	}

// 	if user.SetupToken == nil {
// 		return utils.NewBadRequestError("No setup token available")
// 	}

// 	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
// 	if err != nil {
// 		return utils.NewNotFoundError("Employee not found")
// 	}

// 	// Send account setup email
// 	return s.emailService.SendAccountSetupEmail(user.Email, employee.Name, *user.SetupToken)
// }