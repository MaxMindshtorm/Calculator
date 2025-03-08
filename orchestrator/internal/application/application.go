package application

import "github.com/MaxMindshtorm/calculator/orchestrator/internal/config"

type Application struct {
	Cfg *config.Config
}

func New() *Application {
	return &Application{
		Cfg: config.ConfigFromEnv(),
	}
}
