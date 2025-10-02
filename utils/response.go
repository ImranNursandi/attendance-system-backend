package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessJSON(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorJSON(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}