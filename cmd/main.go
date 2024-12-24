package main

import (
	"fmt"

	"github.com/MaxMindshtorm/calculator/internal/application"
)

func main() {
	a := application.New()
	a.StartServer()
	fmt.Println("Hello world")
}
