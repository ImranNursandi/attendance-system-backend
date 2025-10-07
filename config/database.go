// database/database.go
package database

import (
	"attendance-system/config"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	cfg := config.GetConfig()
	
	var dsn string
	
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		// Add MySQL parameters if not present
		if !strings.Contains(dsn, "?") {
			dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
		} else {
			dsn += "&charset=utf8mb4&parseTime=True&loc=Local"
		}
	} else {
		// Fallback to individual connection parameters
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	}

	log.Printf("🔗 Connecting to database...")
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("✅ Database connected successfully")
	return db, nil
}