package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Response
	Pagination interface{} `json:"pagination,omitempty"`
}

// SuccessJSON sends a successful JSON response
func SuccessJSON(ctx *gin.Context, statusCode int, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	ctx.JSON(statusCode, response)
}

// SuccessPaginatedJSON sends a successful paginated JSON response
func SuccessPaginatedJSON(ctx *gin.Context, statusCode int, message string, data interface{}, pagination interface{}) {
	response := PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Pagination: pagination,
	}
	ctx.JSON(statusCode, response)
}

// ErrorJSON sends an error JSON response
func ErrorJSON(ctx *gin.Context, statusCode int, errorMessage string) {
	response := Response{
		Success: false,
		Message: "Error",
		Error:   errorMessage,
	}
	ctx.JSON(statusCode, response)
}

// ErrorJSONWithDetails sends an error JSON response with additional details
func ErrorJSONWithDetails(ctx *gin.Context, statusCode int, errorMessage string, details interface{}) {
	response := Response{
		Success: false,
		Message: "Error",
		Error:   errorMessage,
		Data:    details,
	}
	ctx.JSON(statusCode, response)
}

// HandleError handles different types of errors and sends appropriate response
func HandleError(ctx *gin.Context, err error) {
	switch e := err.(type) {
	case *CustomError:
		ErrorJSON(ctx, e.StatusCode, e.Message)
	case error:
		// Check if it's a database constraint error
		if IsDuplicateError(err) {
			ErrorJSON(ctx, http.StatusConflict, "Duplicate entry found")
			return
		}
		if IsForeignKeyError(err) {
			ErrorJSON(ctx, http.StatusBadRequest, "Referenced record not found")
			return
		}
		ErrorJSON(ctx, http.StatusInternalServerError, "Internal server error")
	default:
		ErrorJSON(ctx, http.StatusInternalServerError, "Unknown error occurred")
	}
}