package middleware

import (
	"attendance-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateJSON validates the request body against the provided struct
func ValidateJSON(schema interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(schema); err != nil {
			handleValidationError(c, err)
			return
		}
		c.Set("validatedBody", schema)
		c.Next()
	}
}

// ValidateQuery validates query parameters against the provided struct
func ValidateQuery(schema interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(schema); err != nil {
			handleValidationError(c, err)
			return
		}
		c.Set("validatedQuery", schema)
		c.Next()
	}
}

// handleValidationError processes validation errors and returns formatted response
func handleValidationError(c *gin.Context, err error) {
	var validationErrors []ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fe.Field(),
				Message: getValidationMessage(fe),
			})
		}
	} else {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "body",
			Message: err.Error(),
		})
	}

	utils.ErrorJSONWithDetails(c, http.StatusBadRequest, "Validation failed", validationErrors)
	c.Abort()
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "numeric":
		return fe.Field() + " must be a number"
	case "alphanum":
		return fe.Field() + " must contain only letters and numbers"
	default:
		return fe.Field() + " is invalid"
	}
}

// RegisterCustomValidators registers custom validation rules
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators here
		v.RegisterValidation("time_format", validateTimeFormat)
		v.RegisterValidation("employee_id", validateEmployeeID)
	}
}

// validateTimeFormat custom validator for time format (HH:MM:SS)
func validateTimeFormat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) != 8 {
		return false
	}
	if value[2] != ':' || value[5] != ':' {
		return false
	}
	return true
}

// validateEmployeeID custom validator for employee ID format
func validateEmployeeID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) < 3 || len(value) > 50 {
		return false
	}
	// Add more specific validation rules for employee ID format
	return true
}