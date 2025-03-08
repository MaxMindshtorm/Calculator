package main

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/application"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/transport/router"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("Config.env")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	app := application.New()
	router.Run(app)
}
