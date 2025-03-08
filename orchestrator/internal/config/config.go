package config

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/logic"
	"os"
	"strconv"
)

type Config struct {
	OrchestratorPort int
}

func ConfigFromEnv() *Config {
	loadTime()
	port := getEnvInt("ORCHESTRATOR_PORT", 8081)
	return &Config{
		OrchestratorPort: port,
	}
}

func loadTime() {
	logic.OperationTimes["+"] = getEnvInt("TIME_ADDITION_MS", 1000)
	logic.OperationTimes["-"] = getEnvInt("TIME_SUBTRACTION_MS", 1000)
	logic.OperationTimes["*"] = getEnvInt("TIME_MULTIPLICATIONS_MS", 2000)
	logic.OperationTimes["/"] = getEnvInt("TIME_DIVISIONS_MS", 3000)
}

func getEnvInt(key string, defaultVal int) int {
	value, exist := os.LookupEnv(key)
	fmt.Println(exist)
	if exist {
		intVal, _ := strconv.Atoi(value)
		return intVal
	}
	return defaultVal
}
