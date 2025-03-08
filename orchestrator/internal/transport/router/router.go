package router

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/application"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/transport/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run(app *application.Application) {

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", handlers.Calculate).Methods("POST")
	r.HandleFunc("/api/v1/expressions", handlers.GetExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", handlers.GetExpression).Methods("GET")

	r.HandleFunc("/internal/task", handlers.GetTask).Methods("GET")
	r.HandleFunc("/internal/task", handlers.SubmitTaskResult).Methods("POST")

	fmt.Println(fmt.Sprintf("Starting orchestrator server on :%d...", app.Cfg.OrchestratorPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", app.Cfg.OrchestratorPort), r))
}
