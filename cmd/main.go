package main

import (
	"fmt"

	"github.com/MaxMindshtorm/calculator/internal/application"
)

func main() {
	cfg := application.Config{
		Port: 8080,
	}
	a := application.New(cfg)
	a.StartServer()
	fmt.Println("Hello world")
}
