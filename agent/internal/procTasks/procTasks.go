package procTasks

import (
	"fmt"
	"github.com/MaxMindshtorm/calculator/agent/internal/models"
	"strconv"
	"time"
)

func ProcessTask(task *models.Task) (float64, error) {
	var arg1, arg2 float64
	var err error

	arg1, err = strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid arg1: %v", err)
	}

	arg2, err = strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid arg2: %v", err)
	}

	var result float64
	switch task.Operation {
	case "+":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result = arg1 + arg2
	case "-":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result = arg1 - arg2
	case "*":
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result = arg1 * arg2
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result = arg1 / arg2
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}

	return result, nil
}
