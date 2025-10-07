package database

import (
	"attendance-system/config"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GetConfig()
	
	var dsn string
	
	// Priority 1: Use DATABASE_URL (Railway)
	if cfg.DatabaseURL != "" {
		dsn = convertRailwayDSN(cfg.DatabaseURL)
		log.Printf("🔗 Using DATABASE_URL for connection: %s", maskPassword(dsn))
	} else {
		// Priority 2: Use individual connection parameters
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		log.Println("🔗 Using individual DB config for connection")
	}

	// Configure GORM
	gormConfig := &gorm.Config{}
	if os.Getenv("GIN_MODE") == "release" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully")
}

// convertRailwayDSN converts Railway MySQL URL to standard MySQL DSN
func convertRailwayDSN(railwayURL string) string {
	// Remove mysql:// prefix
	cleanURL := strings.Replace(railwayURL, "mysql://", "", 1)
	
	// Split into user:pass@host:port/database
	parts := strings.Split(cleanURL, "@")
	if len(parts) != 2 {
		log.Printf("⚠️ Unexpected DSN format, using as-is: %s", maskPassword(railwayURL))
		return railwayURL
	}
	
	userPass := parts[0]
	hostDB := parts[1]
	
	// Extract user and password
	userParts := strings.Split(userPass, ":")
	user := userParts[0]
	password := userParts[1] // This should exist
	
	// Extract host:port and database
	hostParts := strings.Split(hostDB, "/")
	hostPort := hostParts[0]
	database := "railway"
	if len(hostParts) > 1 {
		database = strings.Split(hostParts[1], "?")[0] // Remove query params
	}
	
	// Construct standard MySQL DSN
	standardDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		hostPort,
		database,
	)
	
	return standardDSN
}

// maskPassword hides password in logs
func maskPassword(dsn string) string {
	if strings.Contains(dsn, ":") && strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, ":")
		if len(parts) >= 2 {
			passwordPart := parts[1]
			if strings.Contains(passwordPart, "@") {
				masked := strings.Split(passwordPart, "@")[0]
				if len(masked) > 2 {
					return strings.Replace(dsn, ":"+masked+"@", ":****@", 1)
				}
			}
		}
	}
	return dsn
}

func CloseDB() {
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}