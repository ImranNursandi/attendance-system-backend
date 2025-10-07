// database/migrate.go
package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	// Read migration file
	migrationSQL, err := os.ReadFile("database/migration.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Split into individual statements
	statements := strings.Split(string(migrationSQL), ";")

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		log.Printf("Executing migration statement %d...", i+1)
		if _, err := sqlDB.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute statement %d: %v\nStatement: %s", i+1, err, statement)
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}