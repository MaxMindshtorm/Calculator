package logic

import (
	"github.com/MaxMindshtorm/calculator/orchestrator/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	OperationTimes = map[string]int{
		"+": 100,
		"-": 100,
		"*": 200,
		"/": 200,
	}
}

func TestCreateTasksFromPostfix(t *testing.T) {
	tests := []struct {
		name           string
		postfix        []string
		expectedTasks  int
		expectedResult *float64
		wantError      bool
	}{
		{
			name:          "Simple addition",
			postfix:       []string{"2", "3", "+"},
			expectedTasks: 1,
			wantError:     false,
		},
		{
			name:          "Multiple operations",
			postfix:       []string{"2", "3", "+", "4", "*"},
			expectedTasks: 2,
			wantError:     false,
		},
		{
			name:           "Single value",
			postfix:        []string{"42"},
			expectedTasks:  0,
			expectedResult: floatPtr(42),
			wantError:      false,
		},
		{
			name:      "Invalid postfix - not enough operands",
			postfix:   []string{"+"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &models.Expression{
				ID:      "test-expr-id",
				Postfix: tt.postfix,
				TaskGraph: &models.TaskGraph{
					Tasks:      make(map[string]*models.Task),
					Results:    make(map[string]float64),
					Ready:      make([]*models.Task, 0),
					InProgress: make(map[string]bool),
					Completed:  make(map[string]bool),
				},
				Tasks: make([]*models.Task, 0),
			}

			err := createTasksFromPostfix(expr)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTasks, len(expr.Tasks))

			if tt.expectedResult != nil {
				assert.NotNil(t, expr.Result)
				assert.Equal(t, *tt.expectedResult, *expr.Result)
				assert.Equal(t, models.StatusCompleted, expr.Status)
			}
		})
	}
}

func TestAddInitialTasksToQueue(t *testing.T) {
	tests := []struct {
		name                string
		setup               func() *models.Expression
		expectedQueueLength int
		expectedStatus      models.ExpressionStatus
	}{
		{
			name: "Add one task with no dependencies",
			setup: func() *models.Expression {
				TaskQueue = []*models.Task{}
				expr := &models.Expression{
					ID:     "test-expr-id",
					Status: models.StatusPending,
					TaskGraph: &models.TaskGraph{
						Tasks:      make(map[string]*models.Task),
						Results:    make(map[string]float64),
						Ready:      make([]*models.Task, 0),
						InProgress: make(map[string]bool),
						Completed:  make(map[string]bool),
					},
				}

				task := &models.Task{
					ID:           "task1",
					ExpressionID: expr.ID,
					Arg1:         "2",
					Arg2:         "3",
					Operation:    "+",
					Status:       "pending",
					Dependencies: []string{},
					ResultVar:    "var0",
				}

				expr.Tasks = append(expr.Tasks, task)
				expr.TaskGraph.Tasks[task.ID] = task

				return expr
			},
			expectedQueueLength: 1,
			expectedStatus:      models.StatusProcessing,
		},
		{
			name: "Expression already completed",
			setup: func() *models.Expression {
				TaskQueue = []*models.Task{}
				result := 5.0
				expr := &models.Expression{
					ID:     "test-expr-id",
					Status: models.StatusCompleted,
					Result: &result,
					TaskGraph: &models.TaskGraph{
						Tasks:      make(map[string]*models.Task),
						Results:    make(map[string]float64),
						Ready:      make([]*models.Task, 0),
						InProgress: make(map[string]bool),
						Completed:  make(map[string]bool),
					},
				}
				return expr
			},
			expectedQueueLength: 0,
			expectedStatus:      models.StatusCompleted,
		},
		{
			name: "Tasks with dependencies",
			setup: func() *models.Expression {
				TaskQueue = []*models.Task{}
				expr := &models.Expression{
					ID:     "test-expr-id",
					Status: models.StatusPending,
					TaskGraph: &models.TaskGraph{
						Tasks:      make(map[string]*models.Task),
						Results:    make(map[string]float64),
						Ready:      make([]*models.Task, 0),
						InProgress: make(map[string]bool),
						Completed:  make(map[string]bool),
					},
				}

				task1 := &models.Task{
					ID:           "task1",
					ExpressionID: expr.ID,
					Arg1:         "2",
					Arg2:         "3",
					Operation:    "+",
					Status:       "pending",
					Dependencies: []string{},
					ResultVar:    "var0",
				}

				task2 := &models.Task{
					ID:           "task2",
					ExpressionID: expr.ID,
					Arg1:         "var0",
					Arg2:         "4",
					Operation:    "*",
					Status:       "pending",
					Dependencies: []string{"task1"},
					ResultVar:    "var1",
				}

				expr.Tasks = append(expr.Tasks, task1, task2)
				expr.TaskGraph.Tasks[task1.ID] = task1
				expr.TaskGraph.Tasks[task2.ID] = task2

				return expr
			},
			expectedQueueLength: 1,
			expectedStatus:      models.StatusProcessing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.setup()
			AddInitialTasksToQueue(expr)

			assert.Equal(t, tt.expectedQueueLength, len(TaskQueue))
			assert.Equal(t, tt.expectedStatus, expr.Status)

			if tt.expectedQueueLength > 0 {
				assert.Equal(t, tt.expectedQueueLength, len(expr.TaskGraph.Ready))
			}
		})
	}
}

func TestUpdateTaskArguments(t *testing.T) {
	tests := []struct {
		name         string
		task         *models.Task
		expr         *models.Expression
		expectedArg1 string
		expectedArg2 string
	}{
		{
			name: "Numeric arguments remain unchanged",
			task: &models.Task{
				ID:           "task1",
				ExpressionID: "expr1",
				Arg1:         "2",
				Arg2:         "3",
				Operation:    "+",
			},
			expr: &models.Expression{
				ID: "expr1",
				TaskGraph: &models.TaskGraph{
					Results: map[string]float64{},
				},
			},
			expectedArg1: "2",
			expectedArg2: "3",
		},
		{
			name: "Variable arguments are replaced with values",
			task: &models.Task{
				ID:           "task1",
				ExpressionID: "expr1",
				Arg1:         "var0",
				Arg2:         "var1",
				Operation:    "+",
			},
			expr: &models.Expression{
				ID: "expr1",
				TaskGraph: &models.TaskGraph{
					Results: map[string]float64{
						"var0": 5.0,
						"var1": 7.0,
					},
				},
			},
			expectedArg1: "5",
			expectedArg2: "7",
		},
		{
			name: "Mix of numeric and variable arguments",
			task: &models.Task{
				ID:           "task1",
				ExpressionID: "expr1",
				Arg1:         "var0",
				Arg2:         "3",
				Operation:    "+",
			},
			expr: &models.Expression{
				ID: "expr1",
				TaskGraph: &models.TaskGraph{
					Results: map[string]float64{
						"var0": 5.0,
					},
				},
			},
			expectedArg1: "5",
			expectedArg2: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateTaskArguments(tt.task, tt.expr)
			assert.Equal(t, tt.expectedArg1, tt.task.Arg1)
			assert.Equal(t, tt.expectedArg2, tt.task.Arg2)
		})
	}
}

func TestAddNextTasksToQueue(t *testing.T) {
	tests := []struct {
		name                string
		setup               func() *models.Expression
		expectedQueueLength int
		expectedReadyLength int
	}{
		{
			name: "Add task when dependencies are satisfied",
			setup: func() *models.Expression {
				TaskQueue = []*models.Task{}
				expr := &models.Expression{
					ID:     "test-expr-id",
					Status: models.StatusProcessing,
					TaskGraph: &models.TaskGraph{
						Tasks:      make(map[string]*models.Task),
						Results:    map[string]float64{"var0": 8.0},
						Ready:      make([]*models.Task, 0),
						InProgress: make(map[string]bool),
						Completed:  map[string]bool{"task1": true},
					},
				}

				task1 := &models.Task{
					ID:           "task1",
					ExpressionID: expr.ID,
					Arg1:         "2",
					Arg2:         "3",
					Operation:    "+",
					Status:       "completed",
					Dependencies: []string{},
					ResultVar:    "var0",
				}

				task2 := &models.Task{
					ID:           "task2",
					ExpressionID: expr.ID,
					Arg1:         "var0",
					Arg2:         "4",
					Operation:    "*",
					Status:       "pending",
					Dependencies: []string{"task1"},
					ResultVar:    "var1",
				}

				expr.Tasks = append(expr.Tasks, task1, task2)
				expr.TaskGraph.Tasks[task1.ID] = task1
				expr.TaskGraph.Tasks[task2.ID] = task2

				return expr
			},
			expectedQueueLength: 1,
			expectedReadyLength: 1,
		},
		{
			name: "Don't add tasks that are already in progress",
			setup: func() *models.Expression {
				TaskQueue = []*models.Task{}
				expr := &models.Expression{
					ID:     "test-expr-id",
					Status: models.StatusProcessing,
					TaskGraph: &models.TaskGraph{
						Tasks:      make(map[string]*models.Task),
						Results:    make(map[string]float64),
						Ready:      make([]*models.Task, 0),
						InProgress: map[string]bool{"task1": true},
						Completed:  make(map[string]bool),
					},
				}

				task1 := &models.Task{
					ID:           "task1",
					ExpressionID: expr.ID,
					Arg1:         "2",
					Arg2:         "3",
					Operation:    "+",
					Status:       "in_progress",
					Dependencies: []string{},
					ResultVar:    "var0",
				}

				expr.Tasks = append(expr.Tasks, task1)
				expr.TaskGraph.Tasks[task1.ID] = task1

				return expr
			},
			expectedQueueLength: 0,
			expectedReadyLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.setup()
			AddNextTasksToQueue(expr)

			assert.Equal(t, tt.expectedQueueLength, len(TaskQueue))
			assert.Equal(t, tt.expectedReadyLength, len(expr.TaskGraph.Ready))
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
