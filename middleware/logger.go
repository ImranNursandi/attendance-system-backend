package middleware

import (
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger middleware configuration
func Logger() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithUTC(true),
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Logger()
		}),
	)
}

// CustomLogger provides enhanced logging
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		duration := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Log the request
		logger := zerolog.New(gin.DefaultWriter).With().
			Timestamp().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", statusCode).
			Dur("duration", duration).
			Str("request_id", c.GetHeader("X-Request-ID")).
			Logger()

		// Log based on status code
		switch {
		case statusCode >= 400 && statusCode < 500:
			logger.Warn().Msg("Client error")
		case statusCode >= 500:
			logger.Error().Msg("Server error")
		default:
			logger.Info().Msg("Request completed")
		}
	}
}