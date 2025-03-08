package main

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/agent/internal/application"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("Config.env")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	app := application.New()
	application.Run(app)
}
