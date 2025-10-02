package main

import (
	"attendance-system/config"
	"attendance-system/routes"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	router.Run(":8080")
}