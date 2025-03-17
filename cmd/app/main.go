package main

import (
	"log"

	"family-flow-app/internal/app"
	"github.com/joho/godotenv"
)

const (
	configPath = "config/config.yaml"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Migrate()
	app.Run(configPath)
}
