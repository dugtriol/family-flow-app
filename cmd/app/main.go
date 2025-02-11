package main

import "family-flow-app/internal/app"

const (
	configPath = "config/config.yaml"
)

func main() {
	app.Run(configPath)
}
