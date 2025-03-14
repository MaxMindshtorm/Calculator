package handlers

import (
	"encoding/json"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/logic"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func Calculate(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if req.Expression == "" {
		http.Error(w, "Expression is required", http.StatusUnprocessableEntity)
		return
	}

	expressionID := uuid.New().String()
	expression := &models.Expression{
		ID:         expressionID,
		Expression: req.Expression,
		Status:     models.StatusPending,
	}

	if err := logic.ParseExpression(expression); err != nil {
		http.Error(w, "Invalid expression: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	logic.ExpressionsMutex.Lock()
	logic.Expressions[expressionID] = expression
	logic.ExpressionsMutex.Unlock()

	logic.AddInitialTasksToQueue(expression)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CalculateResponse{ID: expressionID})
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	logic.ExpressionsMutex.RLock()
	defer logic.ExpressionsMutex.RUnlock()

	expressionsList := make([]*models.Expression, 0, len(logic.Expressions))
	for _, expr := range logic.Expressions {
		expressionsList = append(expressionsList, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ExpressionsResponse{Expressions: expressionsList})
}

func GetExpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	logic.ExpressionsMutex.RLock()
	expr, exists := logic.Expressions[id]
	logic.ExpressionsMutex.RUnlock()

	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ExpressionResponse{Expression: expr})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	logic.TaskQueueMutex.Lock()
	defer logic.TaskQueueMutex.Unlock()

	if len(logic.TaskQueue) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	task := logic.TaskQueue[0]
	logic.TaskQueue = logic.TaskQueue[1:]

	expr, exists := logic.Expressions[task.ExpressionID]
	if exists {
		expr.TaskGraph.Mu.Lock()
		expr.TaskGraph.InProgress[task.ID] = true
		expr.TaskGraph.Mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.TaskResponse{Task: task})
}

func SubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var req models.TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	var foundTask *models.Task
	var foundExpr *models.Expression

	logic.ExpressionsMutex.Lock()
	defer logic.ExpressionsMutex.Unlock()

	for _, expr := range logic.Expressions {
		for _, task := range expr.Tasks {
			if task.ID == req.ID {
				foundTask = task
				foundExpr = expr
				break
			}
		}
		if foundTask != nil {
			break
		}
	}

	if foundTask == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	result := req.Result
	foundTask.Result = &result

	foundExpr.TaskGraph.Mu.Lock()
	foundExpr.TaskGraph.Completed[foundTask.ID] = true
	delete(foundExpr.TaskGraph.InProgress, foundTask.ID)
	foundExpr.TaskGraph.Results[foundTask.ResultVar] = result
	foundExpr.TaskGraph.Mu.Unlock()

	allCompleted := true
	for _, task := range foundExpr.Tasks {
		if task.Result == nil {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		finalResult := *foundExpr.Tasks[len(foundExpr.Tasks)-1].Result
		foundExpr.Result = &finalResult
		foundExpr.Status = models.StatusCompleted
	} else {
		logic.AddNextTasksToQueue(foundExpr)
	}

	w.WriteHeader(http.StatusOK)
}
