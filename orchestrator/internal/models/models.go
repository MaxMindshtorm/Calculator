package models

import "sync"

type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression,omitempty"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	Tasks      []*Task          `json:"-"`
	TaskGraph  *TaskGraph       `json:"-"`
	Postfix    []string         `json:"-"`
}

type Task struct {
	ID            string   `json:"id"`
	ExpressionID  string   `json:"-"`
	Arg1          string   `json:"arg1"`
	Arg2          string   `json:"arg2"`
	Operation     string   `json:"operation"`
	OperationTime int      `json:"operation_time"`
	Result        *float64 `json:"result,omitempty"`
	Status        string   `json:"-"`
	Dependencies  []string `json:"-"`
	ResultVar     string   `json:"-"`
}

type TaskGraph struct {
	Tasks      map[string]*Task
	Results    map[string]float64
	Ready      []*Task
	InProgress map[string]bool
	Completed  map[string]bool
	Mu         sync.Mutex
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID string `json:"id"`
}

type ExpressionsResponse struct {
	Expressions []*Expression `json:"expressions"`
}

type ExpressionResponse struct {
	Expression *Expression `json:"expression"`
}

type TaskResponse struct {
	Task *Task `json:"task"`
}

type TaskResultRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}
