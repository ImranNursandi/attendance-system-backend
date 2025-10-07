// database/database.go
package database

import (
	"attendance-system/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GetConfig()
	
	var dsn string
	var dbName string
	
	// Priority 1: Use DATABASE_URL (Railway)
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		dbName = extractDatabaseName(dsn)
		// Add MySQL parameters if not present
		if !strings.Contains(dsn, "?") {
			dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
		} else {
			if !strings.Contains(dsn, "charset=") {
				dsn += "&charset=utf8mb4"
			}
			if !strings.Contains(dsn, "parseTime=") {
				dsn += "&parseTime=True"
			}
			if !strings.Contains(dsn, "loc=") {
				dsn += "&loc=Local"
			}
		}
		log.Println("🔗 Using DATABASE_URL for connection")
	} else {
		// Priority 2: Use individual connection parameters (local development)
		dbName = cfg.DBName
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			dbName,
		)
		log.Println("🔗 Using individual DB config for connection")
	}

	// First, try to connect without database to create it if needed
	if err := createDatabaseIfNotExists(dsn, dbName); err != nil {
		log.Fatal("❌ Failed to create database:", err)
	}

	// Configure GORM logger based on environment
	gormConfig := &gorm.Config{}
	if os.Getenv("GIN_MODE") == "release" || cfg.GinMode == "release" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connected successfully")
}

// createDatabaseIfNotExists creates the database if it doesn't exist
func createDatabaseIfNotExists(dsn, dbName string) error {
	// Extract base DSN without database name
	baseDSN := strings.Split(dsn, "/")[0] + "/"
	
	// Connect to MySQL without specifying database
	db, err := sql.Open("mysql", baseDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		log.Printf("📁 Database '%s' not found, creating...", dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("✅ Database '%s' created successfully", dbName)
	} else {
		log.Printf("✅ Database '%s' already exists", dbName)
	}

	return nil
}

// extractDatabaseName extracts database name from DSN
func extractDatabaseName(dsn string) string {
	// DSN format: mysql://user:pass@host:port/database
	parts := strings.Split(dsn, "/")
	if len(parts) >= 4 {
		// Get the last part and remove any query parameters
		dbPart := strings.Split(parts[3], "?")[0]
		return dbPart
	}
	return "attendance_system" // fallback
}

func CloseDB() {
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			sqlDB.Close()
			log.Println("✅ Database connection closed")
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}