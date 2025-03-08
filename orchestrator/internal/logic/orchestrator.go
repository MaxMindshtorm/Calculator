package logic

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/models"
	"github.com/MaxMindshtorm/calculator/orchestrator/pkg/Calculate"
	"github.com/google/uuid"
	"strconv"
	"sync"
)

var (
	Expressions      = make(map[string]*models.Expression)
	ExpressionsMutex = sync.RWMutex{}
	TaskQueue        = make([]*models.Task, 0)
	TaskQueueMutex   = sync.Mutex{}
	OperationTimes   = make(map[string]int)
)

func ParseExpression(expr *models.Expression) error {
	postfix, err := Calculate.PrefixToPostfix(expr.Expression)
	if err != nil {
		return fmt.Errorf("error parsing expression: %v", err)
	}

	expr.Postfix = postfix

	expr.TaskGraph = &models.TaskGraph{
		Tasks:      make(map[string]*models.Task),
		Results:    make(map[string]float64),
		Ready:      make([]*models.Task, 0),
		InProgress: make(map[string]bool),
		Completed:  make(map[string]bool),
	}

	return createTasksFromPostfix(expr)
}

func createTasksFromPostfix(expr *models.Expression) error {
	var stack []string
	varCounter := 0

	for i, token := range expr.Postfix {
		if isOperator(token) {
			if len(stack) < 2 {
				return fmt.Errorf("invalid expression: not enough operands for operator %s", token)
			}

			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			taskID := uuid.New().String()
			resultVar := fmt.Sprintf("var%d", varCounter)
			varCounter++

			task := &models.Task{
				ID:            taskID,
				ExpressionID:  expr.ID,
				Arg1:          arg1,
				Arg2:          arg2,
				Operation:     token,
				OperationTime: OperationTimes[token],
				Status:        "pending",
				ResultVar:     resultVar,
			}

			for _, otherTask := range expr.Tasks {
				if otherTask.ResultVar == arg1 || otherTask.ResultVar == arg2 {
					task.Dependencies = append(task.Dependencies, otherTask.ID)
				}
			}

			expr.Tasks = append(expr.Tasks, task)
			expr.TaskGraph.Tasks[taskID] = task

			stack = append(stack, resultVar)
		} else {
			stack = append(stack, token)
		}

		if i == len(expr.Postfix)-1 && len(stack) == 1 {
			if len(expr.Tasks) == 0 {
				result, err := strconv.ParseFloat(stack[0], 64)
				if err != nil {
					return fmt.Errorf("invalid number: %v", err)
				}
				expr.Result = &result
				expr.Status = models.StatusCompleted
			}
		}
	}

	return nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func AddInitialTasksToQueue(expr *models.Expression) {
	TaskQueueMutex.Lock()
	defer TaskQueueMutex.Unlock()

	if expr.Status == models.StatusCompleted {
		return
	}

	for _, task := range expr.Tasks {
		if len(task.Dependencies) == 0 {
			updateTaskArguments(task, expr)

			TaskQueue = append(TaskQueue, task)
			expr.TaskGraph.Ready = append(expr.TaskGraph.Ready, task)
		}
	}

	expr.Status = models.StatusProcessing
}

func updateTaskArguments(task *models.Task, expr *models.Expression) {
	if _, err := strconv.ParseFloat(task.Arg1, 64); err != nil {
		if result, exists := expr.TaskGraph.Results[task.Arg1]; exists {
			task.Arg1 = fmt.Sprintf("%g", result)
		}
	}

	if _, err := strconv.ParseFloat(task.Arg2, 64); err != nil {
		if result, exists := expr.TaskGraph.Results[task.Arg2]; exists {
			task.Arg2 = fmt.Sprintf("%g", result)
		}
	}
}

func AddNextTasksToQueue(expr *models.Expression) {
	TaskQueueMutex.Lock()
	defer TaskQueueMutex.Unlock()

	for _, task := range expr.Tasks {
		if expr.TaskGraph.Completed[task.ID] || expr.TaskGraph.InProgress[task.ID] {
			continue
		}

		allDependenciesSatisfied := true
		for _, depID := range task.Dependencies {
			if !expr.TaskGraph.Completed[depID] {
				allDependenciesSatisfied = false
				break
			}
		}

		if allDependenciesSatisfied {
			updateTaskArguments(task, expr)

			TaskQueue = append(TaskQueue, task)
			expr.TaskGraph.Mu.Lock()
			expr.TaskGraph.Ready = append(expr.TaskGraph.Ready, task)
			expr.TaskGraph.Mu.Unlock()
		}
	}
}
