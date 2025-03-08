package models

import (
	"net/http"
	"time"
)

type Task struct {
	ID            string   `json:"id"`
	Arg1          string   `json:"arg1"`
	Arg2          string   `json:"arg2"`
	Operation     string   `json:"operation"`
	OperationTime int      `json:"operation_time"`
	Result        *float64 `json:"result,omitempty"`
}

type TaskResponse struct {
	Task *Task `json:"task"`
}

type TaskResultRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

var (
	orchestratorURL string
	client          = &http.Client{
		Timeout: 10 * time.Second,
	}
)
