package application

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/agent/internal/config"
	"github.com/MaxMindshtorm/calculator/agent/internal/logic"
	"sync"
)

type App struct {
	Cfg *config.Config
}

func New() *App {
	return &App{
		Cfg: config.New(),
	}
}

func Run(app *App) {
	var wg sync.WaitGroup
	for i := 0; i < app.Cfg.ComputingPower; i++ {
		wg.Add(1)
		go logic.Worker(i, &wg, fmt.Sprintf("http://localhost:%d", app.Cfg.OrchestratorPort))
	}

	wg.Wait()
}
