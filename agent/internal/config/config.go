package config

import (
	"os"
	"strconv"
)

type Config struct {
	ComputingPower   int
	OrchestratorPort int
}

func New() *Config {
	CompPow := getEnvInt("COMPUTING_POWER", 4)
	OrchPort := getEnvInt("ORCHESTRATOR_PORT", 8081)
	return &Config{
		ComputingPower:   CompPow,
		OrchestratorPort: OrchPort,
	}

}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
