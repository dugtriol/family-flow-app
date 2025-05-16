package main

import (
	`log`

	_ "family-flow-app/docs"
	"family-flow-app/internal/app"
	`github.com/joho/godotenv`
)

const (
	configPath = "config/config.yaml"
)

// @title Family Flow App
// @version 1.0
// @description This is a sample server FamilyFlow server.

// @host localhost:8080
// @BasePath /api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Migrate()
	app.Run(configPath)
}
